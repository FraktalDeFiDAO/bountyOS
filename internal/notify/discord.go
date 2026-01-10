package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bountyos-v8/internal/core"
	"bountyos-v8/internal/security"
)

type DiscordNotifier struct {
	webhookURL string
	client     *http.Client
}

func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
		client:     security.SecureHTTPClient(),
	}
}

func (n *DiscordNotifier) Alert(bounty core.Bounty) error {
	if n.webhookURL == "" {
		return nil
	}

	color := 0x10b981 // Green
	if bounty.Score >= 80 {
		color = 0xf43f5e // Red
	} else if bounty.Score >= 50 {
		color = 0xfbbf24 // Yellow
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       "🎯 New Bounty Detected!",
				"description": bounty.Title,
				"url":         bounty.URL,
				"color":       color,
				"fields": []map[string]interface{}{
					{"name": "Platform", "value": bounty.Platform, "inline": true},
					{"name": "Reward", "value": fmt.Sprintf("%s %s", bounty.Reward, bounty.Currency), "inline": true},
					{"name": "Score", "value": fmt.Sprintf("%d", bounty.Score), "inline": true},
					{"name": "Payment", "value": bounty.PaymentType, "inline": true},
				},
				"footer": map[string]interface{}{
					"text": "BountyOS v8: Obsidian Sniper",
				},
				"timestamp": time.Now().Format(time.RFC3339),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", n.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("discord returned status %d", resp.StatusCode)
	}

	return nil
}

func (n *DiscordNotifier) Notify(message string) error {
	if n.webhookURL == "" {
		return nil
	}

	payload := map[string]interface{}{
		"content": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
