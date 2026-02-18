package ui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"bountyos-v8/internal/adapters/storage"
	"bountyos-v8/internal/core"
	"bountyos-v8/internal/resilience"
	"bountyos-v8/internal/security"

	"github.com/gorilla/websocket"
)

type WSClientState int

const (
	WSStateConnecting WSClientState = iota
	WSStateConnected
	WSStateDisconnecting
	WSStateDisconnected
)

type WSClient struct {
	conn            *websocket.Conn
	state           WSClientState
	lastActivity    time.Time
	sendQueue       chan []byte
	done            chan struct{}
	mu              sync.RWMutex
	pendingMessages int
}

func NewWSClient(conn *websocket.Conn) *WSClient {
	return &WSClient{
		conn:         conn,
		state:        WSStateConnecting,
		lastActivity: time.Now(),
		sendQueue:    make(chan []byte, 100),
		done:         make(chan struct{}),
	}
}

func (c *WSClient) State() WSClientState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

func (c *WSClient) SetState(state WSClientState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = state
}

func (c *WSClient) LastActivity() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastActivity
}

func (c *WSClient) UpdateActivity() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastActivity = time.Now()
}

func (c *WSClient) QueueSize() int {
	return len(c.sendQueue)
}

func (c *WSClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.state == WSStateDisconnected {
		return nil
	}

	c.state = WSStateDisconnecting
	close(c.done)
	close(c.sendQueue)
	c.state = WSStateDisconnected
	return c.conn.Close()
}

func (c *WSClient) Send(message []byte) bool {
	select {
	case c.sendQueue <- message:
		c.mu.Lock()
		c.pendingMessages++
		c.mu.Unlock()
		return true
	default:
		return false
	}
}

type WebUI struct {
	storage              *storage.SQLiteStorage
	port                 int
	bountiesLimit        int
	statsLimit           int
	fetchIntervalSeconds int
	staticDir            string
	frontendEnabled      bool
	clientsMu            sync.RWMutex
	clients              map[*websocket.Conn]*WSClient
	server               *http.Server
	wg                   sync.WaitGroup

	wsPingInterval  time.Duration
	wsPongTimeout   time.Duration
	wsClientTimeout time.Duration
	wsMaxQueueSize  int
}

func NewWebUI(storageInst *storage.SQLiteStorage, port int, bountiesLimit int, statsLimit int, fetchIntervalSeconds int, staticDir string) *WebUI {
	if bountiesLimit <= 0 {
		bountiesLimit = 1000
	}
	if statsLimit <= 0 {
		statsLimit = 100
	}
	if fetchIntervalSeconds <= 0 {
		fetchIntervalSeconds = 5
	}

	return &WebUI{
		storage:              storageInst,
		port:                 port,
		bountiesLimit:        bountiesLimit,
		statsLimit:           statsLimit,
		fetchIntervalSeconds: fetchIntervalSeconds,
		staticDir:            staticDir,
		clients:              make(map[*websocket.Conn]*WSClient),
		wsPingInterval:       30 * time.Second,
		wsPongTimeout:        10 * time.Second,
		wsClientTimeout:      60 * time.Second,
		wsMaxQueueSize:       100,
	}
}

func (ui *WebUI) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/bounties", ui.handleBounties)
	mux.HandleFunc("/api/stats", ui.handleStats)
	mux.HandleFunc("/ws", ui.handleWS)
	mux.HandleFunc("/health", ui.handleHealth)
	mux.HandleFunc("/health/ready", ui.handleReady)
	mux.HandleFunc("/health/live", ui.handleLive)
	mux.HandleFunc("/", ui.handleIndex)

	ui.frontendEnabled = ui.resolveStaticDir()

	ui.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", ui.port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	security.GetLogger().Info("Starting Web UI on http://localhost:%d", ui.port)

	go func() {
		if err := ui.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			security.GetLogger().Error("Web UI server error: %v", err)
		}
	}()

	return nil
}

func (ui *WebUI) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ui.closeAllClients()

	if ui.server != nil {
		if err := ui.server.Shutdown(ctx); err != nil {
			security.GetLogger().Error("Error shutting down Web UI: %v", err)
			return err
		}
	}

	ui.wg.Wait()
	return nil
}

func (ui *WebUI) handleBounties(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	bounties, err := ui.storage.GetRecentContext(ctx, ui.bountiesLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sort.Slice(bounties, func(i, j int) bool {
		return bounties[i].Score > bounties[j].Score
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bounties)
}

func (ui *WebUI) handleStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	bounties, err := ui.storage.GetRecentContext(ctx, ui.statsLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := struct {
		TotalCount  int            `json:"total_count"`
		ByPlatform  map[string]int `json:"by_platform"`
		AvgScore    float64        `json:"avg_score"`
		CryptoCount int            `json:"crypto_count"`
	}{
		ByPlatform: make(map[string]int),
	}

	stats.TotalCount = len(bounties)
	var totalScore int
	for _, b := range bounties {
		stats.ByPlatform[b.Platform]++
		totalScore += b.Score
		if b.PaymentType == "crypto" {
			stats.CryptoCount++
		}
	}

	if stats.TotalCount > 0 {
		stats.AvgScore = float64(totalScore) / float64(stats.TotalCount)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (ui *WebUI) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (ui *WebUI) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	ready := true
	checks := make(map[string]interface{})

	if err := ui.storage.Ping(ctx); err != nil {
		ready = false
		checks["database"] = map[string]interface{}{
			"status": "unhealthy",
			"error":  err.Error(),
		}
	} else {
		checks["database"] = map[string]string{"status": "healthy"}
	}

	openConns, idleConns := ui.storage.Stats()
	checks["database_connections"] = map[string]int{
		"open": openConns,
		"idle": idleConns,
	}

	cbStats := resilience.AllStatsFromRegistry()
	if len(cbStats) > 0 {
		checks["circuit_breakers"] = cbStats
		for name, stats := range cbStats {
			if stats.State == resilience.StateOpen {
				ready = false
				security.GetLogger().Warn("Circuit breaker %s is open", name)
			}
		}
	}

	ui.clientsMu.RLock()
	wsClientCount := len(ui.clients)
	ui.clientsMu.RUnlock()
	checks["websocket_clients"] = wsClientCount

	status := http.StatusOK
	if !ready {
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ready":   ready,
		"checks":  checks,
		"version": "v8",
	})
}

func (ui *WebUI) handleLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "alive",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func (ui *WebUI) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		if ui.frontendEnabled {
			ui.serveStatic(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	if ui.frontendEnabled {
		ui.serveStatic(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>BountyOS v8: Obsidian</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #0f172a; color: #e2e8f0; margin: 0; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        header { display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid #334155; padding-bottom: 20px; margin-bottom: 20px; }
        h1 { color: #10b981; margin: 0; font-size: 24px; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: #1e293b; padding: 20px; border-radius: 8px; border: 1px solid #334155; text-align: center; }
        .stat-value { font-size: 28px; font-weight: bold; color: #38bdf8; }
        .stat-label { font-size: 14px; color: #94a3b8; text-transform: uppercase; margin-top: 5px; }
        table { width: 100%; border-collapse: collapse; background: #1e293b; border-radius: 8px; overflow: hidden; }
        th { background: #334155; text-align: left; padding: 12px 15px; font-size: 14px; text-transform: uppercase; color: #94a3b8; }
        td { padding: 12px 15px; border-bottom: 1px solid #334155; }
        tr:hover { background: #2d3748; }
        .score { font-weight: bold; }
        .score-high { color: #f43f5e; }
        .score-med { color: #fbbf24; }
        .score-low { color: #10b981; }
        .platform { color: #94a3b8; font-size: 12px; }
        .payout { color: #38bdf8; font-weight: bold; }
        .link { color: #6366f1; text-decoration: none; font-size: 12px; }
        .link:hover { text-decoration: underline; }
        .badge { display: inline-block; padding: 2px 8px; border-radius: 4px; font-size: 11px; margin-right: 5px; background: #475569; }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🕷️ BOUNTY OS v8: OBSIDIAN</h1>
            <div id="last-updated" style="color: #64748b; font-size: 12px;"></div>
        </header>

        <div class="stats" id="stats">
            <!-- Stats will be loaded here -->
        </div>

        <table>
            <thead>
                <tr>
                    <th>Score</th>
                    <th>Platform</th>
                    <th>Payout</th>
                    <th>Task</th>
                </tr>
            </thead>
            <tbody id="bounties-body">
                <!-- Bounties will be loaded here -->
            </tbody>
        </table>
    </div>

    <script>
        async function fetchData() {
            try {
                const [bountiesResp, statsResp] = await Promise.all([
                    fetch('/api/bounties'),
                    fetch('/api/stats')
                ]);
                
                const bounties = await bountiesResp.json();
                const stats = await statsResp.json();
                
                updateStats(stats);
                updateBounties(bounties);
                document.getElementById('last-updated').textContent = 'Last updated: ' + new Date().toLocaleTimeString();
            } catch (err) {
                console.error('Error fetching data:', err);
            }
        }

        function updateStats(stats) {
            const statsContainer = document.getElementById('stats');
            statsContainer.innerHTML = ' \
                <div class="stat-card"> \
                    <div class="stat-value">' + stats.total_count + '</div> \
                    <div class="stat-label">Total Bounties</div> \
                </div> \
                <div class="stat-card"> \
                    <div class="stat-value">' + stats.crypto_count + '</div> \
                    <div class="stat-label">Crypto Bounties</div> \
                </div> \
                <div class="stat-card"> \
                    <div class="stat-value">' + stats.avg_score.toFixed(1) + '</div> \
                    <div class="stat-label">Avg Urgency</div> \
                </div> \
                <div class="stat-card"> \
                    <div class="stat-value">' + Object.keys(stats.by_platform).length + '</div> \
                    <div class="stat-label">Sources</div> \
                </div> \
            ';
        }

        function updateBounties(bounties) {
            const tbody = document.getElementById('bounties-body');
            tbody.innerHTML = bounties.map(b => {
                let scoreClass = 'score-low';
                if (b.score >= 80) scoreClass = 'score-high';
                else if (b.score >= 50) scoreClass = 'score-med';
                
                return ' \
                    <tr> \
                        <td><span class="score ' + scoreClass + '">' + b.score + '</span></td> \
                        <td><span class="platform">' + b.platform + '</span></td> \
                        <td><span class="payout">' + b.reward + (b.currency ? ' ' + b.currency : '') + '</span></td> \
                        <td> \
                            <div>' + b.title + '</div> \
                            <a href="' + b.url + '" class="link" target="_blank">' + b.url + '</a> \
                        </td> \
                    </tr> \
                ';
            }).join('');
        }

        fetchData();
        setInterval(fetchData, __FETCH_INTERVAL__);
    </script>
</body>
</html>
`
	html = strings.ReplaceAll(html, "__FETCH_INTERVAL__", strconv.Itoa(ui.fetchIntervalSeconds*1000))
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func (ui *WebUI) resolveStaticDir() bool {
	if strings.TrimSpace(ui.staticDir) == "" {
		return false
	}

	indexPath := filepath.Join(ui.staticDir, "index.html")
	if info, err := os.Stat(indexPath); err == nil && !info.IsDir() {
		return true
	}

	return false
}

func (ui *WebUI) serveStatic(w http.ResponseWriter, r *http.Request) {
	requested := filepath.Clean(r.URL.Path)
	if requested == "." || requested == "/" {
		requested = "/index.html"
	}

	fullPath := filepath.Join(ui.staticDir, requested)
	if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
		http.ServeFile(w, r, fullPath)
		return
	}

	http.ServeFile(w, r, filepath.Join(ui.staticDir, "index.html"))
}

func (ui *WebUI) Broadcast(bounty core.Bounty) {
	payload, err := json.Marshal(struct {
		Type string      `json:"type"`
		Data core.Bounty `json:"data"`
	}{
		Type: "bounty",
		Data: bounty,
	})
	if err != nil {
		security.GetLogger().Warn("Failed to marshal bounty for ws: %v", err)
		return
	}

	clients := ui.snapshotClients()
	for _, client := range clients {
		if client.QueueSize() >= ui.wsMaxQueueSize {
			security.GetLogger().Warn("Client send queue full, dropping client")
			go ui.removeClient(client.conn)
			continue
		}

		if !client.Send(payload) {
			security.GetLogger().Warn("Failed to queue message for client")
			go ui.removeClient(client.conn)
		}
	}
}

var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ui *WebUI) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		security.GetLogger().Warn("WebSocket upgrade failed: %v", err)
		return
	}

	client := NewWSClient(conn)
	ui.addClient(conn, client)
	client.SetState(WSStateConnected)

	ui.wg.Add(2)
	go ui.wsReadPump(client)
	go ui.wsWritePump(client)
}

func (ui *WebUI) wsReadPump(client *WSClient) {
	defer ui.wg.Done()
	defer ui.removeClient(client.conn)

	client.conn.SetReadLimit(512)
	client.conn.SetReadDeadline(time.Now().Add(ui.wsPongTimeout + ui.wsClientTimeout))
	client.conn.SetPongHandler(func(string) error {
		client.UpdateActivity()
		client.conn.SetReadDeadline(time.Now().Add(ui.wsPongTimeout + ui.wsClientTimeout))
		return nil
	})

	for {
		_, _, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				security.GetLogger().Debug("WebSocket read error: %v", err)
			}
			return
		}
		client.UpdateActivity()
	}
}

func (ui *WebUI) wsWritePump(client *WSClient) {
	defer ui.wg.Done()
	defer ui.removeClient(client.conn)

	ticker := time.NewTicker(ui.wsPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-client.done:
			client.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case msg, ok := <-client.sendQueue:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				security.GetLogger().Debug("WebSocket write error: %v", err)
				return
			}

			client.mu.Lock()
			if client.pendingMessages > 0 {
				client.pendingMessages--
			}
			client.mu.Unlock()

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (ui *WebUI) addClient(conn *websocket.Conn, client *WSClient) {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()
	ui.clients[conn] = client
}

func (ui *WebUI) removeClient(conn *websocket.Conn) {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()

	if client, ok := ui.clients[conn]; ok {
		client.Close()
		delete(ui.clients, conn)
	}
}

func (ui *WebUI) closeAllClients() {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()

	for conn, client := range ui.clients {
		client.Close()
		delete(ui.clients, conn)
	}
}

func (ui *WebUI) snapshotClients() []*WSClient {
	ui.clientsMu.RLock()
	defer ui.clientsMu.RUnlock()
	out := make([]*WSClient, 0, len(ui.clients))
	for _, client := range ui.clients {
		out = append(out, client)
	}
	return out
}
