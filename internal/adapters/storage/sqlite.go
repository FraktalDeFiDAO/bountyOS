package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"bountyos-v8/internal/core"
	"bountyos-v8/internal/resilience"
	"bountyos-v8/internal/security"

	_ "github.com/mattn/go-sqlite3"
)

var (
	errDBClosed = errors.New("database is closed")
)

type DBConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
	RequestTimeout  time.Duration
}

func DefaultDBConfig() DBConfig {
	return DBConfig{
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		RequestTimeout:  5 * time.Second,
	}
}

type SQLiteStorage struct {
	db          *sql.DB
	config      DBConfig
	stmtCache   map[string]*sql.Stmt
	stmtMu      sync.RWMutex
	closed      bool
	closeMu     sync.RWMutex
	healthCheck *resilience.CircuitBreaker
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	return NewSQLiteStorageWithConfig(dbPath, DefaultDBConfig())
}

func NewSQLiteStorageWithConfig(dbPath string, config DBConfig) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}

	query := `CREATE TABLE IF NOT EXISTS bounties (
		url TEXT PRIMARY KEY,
		title TEXT,
		platform TEXT,
		reward TEXT,
		currency TEXT,
		created_at DATETIME,
		score INTEGER,
		description TEXT,
		tags TEXT,
		expires_at DATETIME,
		payment_type TEXT
	);`

	if _, err := db.ExecContext(ctx, query); err != nil {
		db.Close()
		return nil, err
	}

	return &SQLiteStorage{
		db:          db,
		config:      config,
		stmtCache:   make(map[string]*sql.Stmt),
		healthCheck: resilience.NewCircuitBreaker(3, 1, 10*time.Second),
	}, nil
}

func (s *SQLiteStorage) prepareStatement(ctx context.Context, query string) (*sql.Stmt, error) {
	s.stmtMu.RLock()
	stmt, exists := s.stmtCache[query]
	s.stmtMu.RUnlock()

	if exists {
		return stmt, nil
	}

	s.stmtMu.Lock()
	defer s.stmtMu.Unlock()

	if stmt, exists = s.stmtCache[query]; exists {
		return stmt, nil
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	s.stmtCache[query] = stmt
	return stmt, nil
}

func (s *SQLiteStorage) isClosed() bool {
	s.closeMu.RLock()
	defer s.closeMu.RUnlock()
	return s.closed
}

func (s *SQLiteStorage) Ping(ctx context.Context) error {
	if s.isClosed() {
		return errDBClosed
	}

	return s.healthCheck.Execute(ctx, func() error {
		ctx, cancel := context.WithTimeout(ctx, s.config.RequestTimeout)
		defer cancel()
		return s.db.PingContext(ctx)
	})
}

func (s *SQLiteStorage) Close() error {
	s.closeMu.Lock()
	defer s.closeMu.Unlock()

	if s.closed {
		return nil
	}
	s.closed = true

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		s.stmtMu.Lock()
		defer s.stmtMu.Unlock()

		for _, stmt := range s.stmtCache {
			stmt.Close()
		}
		s.stmtCache = make(map[string]*sql.Stmt)

		done <- s.db.Close()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.New("close timeout exceeded")
	}
}

func (s *SQLiteStorage) withRetry(ctx context.Context, fn func(ctx context.Context) error) error {
	if s.isClosed() {
		return errDBClosed
	}

	cfg := resilience.DefaultRetryConfig()
	cfg.MaxRetries = 2
	cfg.RetryableError = func(err error) bool {
		if errors.Is(err, sql.ErrConnDone) || errors.Is(err, sql.ErrTxDone) {
			return true
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return true
		}
		return false
	}

	return resilience.Retry(ctx, cfg, func() error {
		ctx, cancel := context.WithTimeout(ctx, s.config.RequestTimeout)
		defer cancel()
		return fn(ctx)
	})
}

func (s *SQLiteStorage) Save(bounty core.Bounty) error {
	return s.SaveContext(context.Background(), bounty)
}

func (s *SQLiteStorage) SaveContext(ctx context.Context, bounty core.Bounty) error {
	tagsJSON, err := json.Marshal(bounty.Tags)
	if err != nil {
		return err
	}

	query := `INSERT OR REPLACE INTO bounties 
		(url, title, platform, reward, currency, created_at, score, description, tags, expires_at, payment_type) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	var expiresAt *string
	if bounty.ExpiresAt != nil {
		expireStr := bounty.ExpiresAt.Format(time.RFC3339)
		expiresAt = &expireStr
	}

	return s.withRetry(ctx, func(ctx context.Context) error {
		stmt, err := s.prepareStatement(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx,
			bounty.URL,
			bounty.Title,
			bounty.Platform,
			bounty.Reward,
			bounty.Currency,
			bounty.CreatedAt.Format(time.RFC3339),
			bounty.Score,
			bounty.Description,
			string(tagsJSON),
			expiresAt,
			bounty.PaymentType,
		)
		return err
	})
}

func (s *SQLiteStorage) IsNew(url string) (bool, error) {
	return s.IsNewContext(context.Background(), url)
}

func (s *SQLiteStorage) IsNewContext(ctx context.Context, url string) (bool, error) {
	var exists int

	err := s.withRetry(ctx, func(ctx context.Context) error {
		query := "SELECT 1 FROM bounties WHERE url = ?"
		stmt, err := s.prepareStatement(ctx, query)
		if err != nil {
			return err
		}
		return stmt.QueryRowContext(ctx, url).Scan(&exists)
	})

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists == 0, nil
}

func (s *SQLiteStorage) GetRecent(limit int) ([]core.Bounty, error) {
	return s.GetRecentContext(context.Background(), limit)
}

func (s *SQLiteStorage) GetRecentContext(ctx context.Context, limit int) ([]core.Bounty, error) {
	query := `SELECT url, title, platform, reward, currency, created_at, score, description, tags, expires_at, payment_type
		FROM bounties 
		ORDER BY created_at DESC 
		LIMIT ?`

	var bounties []core.Bounty

	err := s.withRetry(ctx, func(ctx context.Context) error {
		stmt, err := s.prepareStatement(ctx, query)
		if err != nil {
			return err
		}

		rows, err := stmt.QueryContext(ctx, limit)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var bounty core.Bounty
			var createdAtStr, expiresAtStr sql.NullString
			var tagsStr sql.NullString

			err := rows.Scan(
				&bounty.URL,
				&bounty.Title,
				&bounty.Platform,
				&bounty.Reward,
				&bounty.Currency,
				&createdAtStr,
				&bounty.Score,
				&bounty.Description,
				&tagsStr,
				&expiresAtStr,
				&bounty.PaymentType,
			)
			if err != nil {
				security.GetLogger().Error("Error scanning bounty: %v", err)
				continue
			}

			if createdAtStr.Valid {
				bounty.CreatedAt, err = parseTime(createdAtStr.String)
				if err != nil {
					security.GetLogger().Error("Error parsing created_at: %v", err)
					continue
				}
			}

			if expiresAtStr.Valid {
				expireTime, err := parseTime(expiresAtStr.String)
				if err == nil {
					bounty.ExpiresAt = &expireTime
				}
			}

			if tagsStr.Valid {
				var tags []string
				err := json.Unmarshal([]byte(tagsStr.String), &tags)
				if err == nil {
					bounty.Tags = tags
				}
			}

			bounties = append(bounties, bounty)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, err
	}

	return bounties, nil
}

func (s *SQLiteStorage) PurgeInvalidURLs(ctx context.Context, validateHTTP bool, timeout time.Duration) (int, error) {
	if s.isClosed() {
		return 0, errDBClosed
	}

	rows, err := s.db.QueryContext(ctx, "SELECT url FROM bounties")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	removed := 0
	for rows.Next() {
		var urlStr string
		if err := rows.Scan(&urlStr); err != nil {
			continue
		}

		normalized := security.NormalizeURL(urlStr)
		if normalized == "" {
			if _, err := s.db.ExecContext(ctx, "DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
				removed++
			}
			continue
		}

		if normalized != urlStr {
			if _, err := s.db.ExecContext(ctx, "UPDATE bounties SET url = ? WHERE url = ?", normalized, urlStr); err == nil {
				urlStr = normalized
			}
		}

		if !security.ValidateURL(urlStr) {
			if _, err := s.db.ExecContext(ctx, "DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
				removed++
			}
			continue
		}

		if validateHTTP {
			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			ok := security.ValidateURLReachable(checkCtx, urlStr, timeout)
			cancel()
			if !ok {
				if _, err := s.db.ExecContext(ctx, "DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
					removed++
				}
			}
		}
	}

	return removed, nil
}

func (s *SQLiteStorage) HealthCheck(ctx context.Context) error {
	return s.Ping(ctx)
}

func (s *SQLiteStorage) Stats() (openConns, idleConns int) {
	stats := s.db.Stats()
	return stats.OpenConnections, stats.Idle
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
