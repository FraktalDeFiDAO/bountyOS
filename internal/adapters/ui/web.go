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
	"bountyos-v8/internal/security"

	"github.com/gorilla/websocket"
)

type WebUI struct {
	storage              *storage.SQLiteStorage
	port                 int
	bountiesLimit        int
	statsLimit           int
	fetchIntervalSeconds int
	staticDir            string
	frontendEnabled      bool
	clientsMu            sync.Mutex
	clients              map[*websocket.Conn]struct{}
	server               *http.Server
}

func NewWebUI(storage *storage.SQLiteStorage, port int, bountiesLimit int, statsLimit int, fetchIntervalSeconds int, staticDir string) *WebUI {
	if bountiesLimit <= 0 {
		bountiesLimit = 50
	}
	if statsLimit <= 0 {
		statsLimit = 100
	}
	if fetchIntervalSeconds <= 0 {
		fetchIntervalSeconds = 5
	}

	return &WebUI{
		storage:              storage,
		port:                 port,
		bountiesLimit:        bountiesLimit,
		statsLimit:           statsLimit,
		fetchIntervalSeconds: fetchIntervalSeconds,
		staticDir:            staticDir,
		clients:              make(map[*websocket.Conn]struct{}),
	}
}

func (ui *WebUI) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/bounties", ui.handleBounties)
	mux.HandleFunc("/api/stats", ui.handleStats)
	mux.HandleFunc("/ws", ui.handleWS)

	// Static files (placeholder for now)
	mux.HandleFunc("/", ui.handleIndex)

	ui.frontendEnabled = ui.resolveStaticDir()

	ui.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", ui.port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
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
	if ui.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return ui.server.Shutdown(ctx)
	}
	return nil
}

func (ui *WebUI) handleBounties(w http.ResponseWriter, r *http.Request) {
	bounties, err := ui.storage.GetRecent(ui.bountiesLimit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Sort by score
	sort.Slice(bounties, func(i, j int) bool {
		return bounties[i].Score > bounties[j].Score
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bounties)
}

func (ui *WebUI) handleStats(w http.ResponseWriter, r *http.Request) {
	bounties, err := ui.storage.GetRecent(ui.statsLimit)
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

	// Simple embedded HTML fallback
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

	// SPA fallback
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
	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
			ui.removeClient(conn)
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

	ui.addClient(conn)
	defer ui.removeClient(conn)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (ui *WebUI) addClient(conn *websocket.Conn) {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()
	ui.clients[conn] = struct{}{}
}

func (ui *WebUI) removeClient(conn *websocket.Conn) {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()
	if _, ok := ui.clients[conn]; ok {
		delete(ui.clients, conn)
	}
	_ = conn.Close()
}

func (ui *WebUI) snapshotClients() []*websocket.Conn {
	ui.clientsMu.Lock()
	defer ui.clientsMu.Unlock()
	out := make([]*websocket.Conn, 0, len(ui.clients))
	for conn := range ui.clients {
		out = append(out, conn)
	}
	return out
}
