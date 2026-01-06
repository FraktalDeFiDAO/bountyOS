package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// --- CONFIGURATION ---
const (
	PollInterval = 60 * time.Second
	// Add your GitHub token if you hit rate limits: "token YOUR_PAT_HERE"
	GithubToken = ""
)

// --- DATA MODELS ---
type Bounty struct {
	Source   string
	Title    string
	URL      string
	Reward   string
	Platform string
	Urgency  string // "HIGH" or "NORMAL"
}

// --- LOGGING SETUP ---
func setupLogging() *log.Logger {
	file, err := os.OpenFile("bounty_sniper.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	logger := log.New(file, "BOUNTY-SNIPER: ", log.LstdFlags|log.Lshortfile)
	return logger
}

// --- MODULE 1: GITHUB SNIPER ---
func scanGitHub(ch chan<- Bounty, logger *log.Logger) {
	logger.Println("Starting GitHub scan...")

	// We search for issues with 'bounty' label created recently
	// We also look for 'USDC' or '$' in the title to ensure it's funded
	query := "is:issue is:open label:bounty sort:created-desc"
	url := fmt.Sprintf("https://api.github.com/search/issues?q=%s", strings.ReplaceAll(query, " ", "+"))

	req, _ := http.NewRequest("GET", url, nil)
	if GithubToken != "" {
		req.Header.Set("Authorization", GithubToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("[GitHub] Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Items []struct {
			Title     string `json:"title"`
			HTMLURL   string `json:"html_url"`
			CreatedAt string `json:"created_at"`
			Body      string `json:"body"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Printf("[GitHub] JSON decode error: %v\n", err)
		return
	}

	logger.Printf("[GitHub] Found %d items\n", len(result.Items))

	for _, item := range result.Items {
		// Filter: Only items from the last 24 hours
		t, _ := time.Parse(time.RFC3339, item.CreatedAt)
		if time.Since(t) < 24*time.Hour {
			logger.Printf("[GitHub] Processing item: %s\n", item.Title)

			urgency := "NORMAL"
			lowerTitle := strings.ToLower(item.Title)
			if strings.Contains(lowerTitle, "urgent") || strings.Contains(lowerTitle, "fix") {
				urgency = "HIGH"
			}

			ch <- Bounty{
				Source:   "GITHUB",
				Title:    item.Title,
				URL:      item.HTMLURL,
				Reward:   "Unknown (Check Labels)",
				Platform: "GitHub",
				Urgency:  urgency,
			}
		}
	}
	logger.Println("Completed GitHub scan")
}

// --- MODULE 2: SUPERTEAM SNIPER ---
func scanSuperteam(ch chan<- Bounty, logger *log.Logger) {
	logger.Println("Starting Superteam scan...")

	// Undocumented API endpoint for Superteam listings
	url := "https://earn.superteam.fun/api/bounties?active=true&take=20"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logger.Printf("[Superteam] Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check if response is HTML (indicating an error page)
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		logger.Printf("[Superteam] Received HTML response instead of JSON, possibly rate limited or API changed\n")
		return
	}

	// Generic struct to handle their JSON
	var result struct {
		Bounties []struct {
			Title        string      `json:"title"`
			Slug         string      `json:"slug"`
			RewardAmount interface{} `json:"rewardAmount"`
			Token        string      `json:"token"`
			Type         string      `json:"type"` // "bounty" or "project"
			Deadline     string      `json:"deadline"`
		} `json:"bounties"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("[Superteam] Error reading response body: %v\n", err)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Printf("[Superteam] JSON unmarshal error: %v\n", err)
		logger.Printf("[Superteam] Response body: %s\n", string(body))
		return
	}

	logger.Printf("[Superteam] Found %d bounties\n", len(result.Bounties))

	for _, item := range result.Bounties {
		// Filter: We only want "Bounties" (Short term), not "Projects" (Long term)
		if item.Type == "bounty" {
			logger.Printf("[Superteam] Processing bounty: %s\n", item.Title)

			ch <- Bounty{
				Source:   "SUPERTEAM",
				Title:    item.Title,
				URL:      "https://earn.superteam.fun/listings/bounty/" + item.Slug,
				Reward:   fmt.Sprintf("%v %s", item.RewardAmount, item.Token),
				Platform: "Solana",
				Urgency:  "NORMAL", // Can add deadline logic here
			}
		}
	}
	logger.Println("Completed Superteam scan")
}

// --- MODULE 3: ALGORA SNIPER ---
func scanAlgora(ch chan<- Bounty, logger *log.Logger) {
	logger.Println("Starting Algora scan...")

	// Algora API endpoint
	url := "https://api.algora.io/graphql"

	// GraphQL query to get open bounties
	query := `
	{
		bounties(filters: { status: Open }) {
			nodes {
				id
				title
				description
				reward {
					amount
					token
				}
				createdAt
				url
			}
		}
	}`

	requestBody := map[string]string{
		"query": query,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Printf("[Algora] Error marshaling request: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		logger.Printf("[Algora] Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Printf("[Algora] Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Check if response is HTML (indicating an error page)
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") {
		logger.Printf("[Algora] Received HTML response instead of JSON, possibly rate limited or API changed\n")
		return
	}

	var result struct {
		Data struct {
			Bounties struct {
				Nodes []struct {
					ID          string `json:"id"`
					Title       string `json:"title"`
					Description string `json:"description"`
					Reward      struct {
						Amount string `json:"amount"`
						Token  string `json:"token"`
					} `json:"reward"`
					CreatedAt string `json:"createdAt"`
					URL       string `json:"url"`
				} `json:"nodes"`
			} `json:"bounties"`
		} `json:"data"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("[Algora] Error reading response body: %v\n", err)
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		logger.Printf("[Algora] JSON unmarshal error: %v\n", err)
		logger.Printf("[Algora] Response body: %s\n", string(body))
		return
	}

	logger.Printf("[Algora] Found %d bounties\n", len(result.Data.Bounties.Nodes))

	for _, item := range result.Data.Bounties.Nodes {
		// Filter: Only items from the last 24 hours
		t, err := time.Parse(time.RFC3339, item.CreatedAt)
		if err != nil {
			continue // Skip if we can't parse the date
		}

		if time.Since(t) < 24*time.Hour {
			logger.Printf("[Algora] Processing bounty: %s\n", item.Title)

			urgency := "NORMAL"
			lowerTitle := strings.ToLower(item.Title)
			if strings.Contains(lowerTitle, "urgent") || strings.Contains(lowerTitle, "fix") {
				urgency = "HIGH"
			}

			ch <- Bounty{
				Source:   "ALGORA",
				Title:    item.Title,
				URL:      item.URL,
				Reward:   fmt.Sprintf("%s %s", item.Reward.Amount, item.Reward.Token),
				Platform: "Algora",
				Urgency:  urgency,
			}
		}
	}
	logger.Println("Completed Algora scan")
}

// --- MAIN ENGINE ---
func main() {
	logger := setupLogging()
	logger.Println("BountyOS Sniper Engine started")

	fmt.Println("\033[32m[SYSTEM ONLINE] BountyOS Sniper Engine v1.0\033[0m")
	fmt.Println("Targeting: GitHub Issues (<24h), Superteam Earn, Algora...")
	fmt.Println("---------------------------------------------------")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	bountyChan := make(chan Bounty)
	seen := make(map[string]bool)

	// Start Pollers
	go func() {
		for {
			select {
			case <-sigChan:
				logger.Println("Shutdown signal received, stopping pollers...")
				return
			default:
				go scanGitHub(bountyChan, logger)
				go scanSuperteam(bountyChan, logger)
				go scanAlgora(bountyChan, logger)
				time.Sleep(PollInterval)
			}
		}
	}()

	// Listen for Hits
	for {
		select {
		case b := <-bountyChan:
			if !seen[b.URL] {
				seen[b.URL] = true
				logger.Printf("New bounty found: %s - %s\n", b.Source, b.Title)

				// KEYWORD FILTER: Only show things relevant to Devs
				keywords := []string{"fix", "bug", "script", "bot", "api", "python", "go", "react", "vue", "integrate", "frontend", "backend", "smart contract", "blockchain"}
				relevant := false

				titleLower := strings.ToLower(b.Title)
				for _, k := range keywords {
					if strings.Contains(titleLower, k) {
						relevant = true
						break
					}
				}

				if relevant {
					// Visual Alert
					color := "\033[36m" // Cyan
					if b.Urgency == "HIGH" {
						color = "\033[31m" // Red for Urgent
					}

					fmt.Printf("%s[%s] %s\033[0m\n", color, b.Source, b.Title)
					fmt.Printf("   💰 Pay: %s\n", b.Reward)
					fmt.Printf("   🔗 Link: %s\n", b.URL)
					fmt.Print("\a") // Audio Beep

					logger.Printf("Alerted user to bounty: %s - %s\n", b.Source, b.Title)
				}
			}
		case <-sigChan:
			logger.Println("Shutdown signal received, exiting...")
			fmt.Println("\n\033[33m[SHUTDOWN] BountyOS Sniper Engine shutting down gracefully...\033[0m")
			return
		}
	}
}