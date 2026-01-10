package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"bountyos-v8/internal/core"
	"bountyos-v8/internal/security"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the table
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
	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Save(bounty core.Bounty) error {
	// Convert tags to JSON string
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
	} else {
		expiresAt = nil
	}

	_, err = s.db.Exec(query,
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
}

func (s *SQLiteStorage) IsNew(url string) (bool, error) {
	var exists int
	err := s.db.QueryRow("SELECT 1 FROM bounties WHERE url = ?", url).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists == 0, nil
}

func (s *SQLiteStorage) GetRecent(limit int) ([]core.Bounty, error) {
	query := `SELECT url, title, platform, reward, currency, created_at, score, description, tags, expires_at, payment_type
		FROM bounties 
		ORDER BY created_at DESC 
		LIMIT ?`

	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bounties []core.Bounty
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

		// Parse created_at
		if createdAtStr.Valid {
			bounty.CreatedAt, err = parseTime(createdAtStr.String)
			if err != nil {
				security.GetLogger().Error("Error parsing created_at: %v", err)
				continue
			}
		}

		// Parse expires_at
		if expiresAtStr.Valid {
			expireTime, err := parseTime(expiresAtStr.String)
			if err == nil {
				bounty.ExpiresAt = &expireTime
			}
		}

		// Parse tags
		if tagsStr.Valid {
			var tags []string
			err := json.Unmarshal([]byte(tagsStr.String), &tags)
			if err == nil {
				bounty.Tags = tags
			}
		}

		bounties = append(bounties, bounty)
	}

	return bounties, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLiteStorage) PurgeInvalidURLs(ctx context.Context, validateHTTP bool, timeout time.Duration) (int, error) {
	rows, err := s.db.Query("SELECT url FROM bounties")
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
			if _, err := s.db.Exec("DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
				removed++
			}
			continue
		}

		if normalized != urlStr {
			if _, err := s.db.Exec("UPDATE bounties SET url = ? WHERE url = ?", normalized, urlStr); err == nil {
				urlStr = normalized
			}
		}

		if !security.ValidateURL(urlStr) {
			if _, err := s.db.Exec("DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
				removed++
			}
			continue
		}

		if validateHTTP {
			checkCtx, cancel := context.WithTimeout(ctx, timeout)
			ok := security.ValidateURLReachable(checkCtx, urlStr, timeout)
			cancel()
			if !ok {
				if _, err := s.db.Exec("DELETE FROM bounties WHERE url = ?", urlStr); err == nil {
					removed++
				}
			}
		}
	}

	return removed, nil
}

func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
