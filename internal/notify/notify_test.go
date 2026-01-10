package notify

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"bountyos-v8/internal/core"
)

func TestDiscordNotifier(t *testing.T) {
	// Mock Discord server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	notifier := NewDiscordNotifier(ts.URL)
	bounty := core.Bounty{
		Title:    "Test Bounty",
		Platform: "GitHub",
		Reward:   "100",
		Currency: "USDC",
		URL:      "https://github.com/test",
		Score:    85,
	}

	err := notifier.Alert(bounty)
	if err != nil {
		t.Errorf("Discord Alert failed: %v", err)
	}

	err = notifier.Notify("Test Message")
	if err != nil {
		t.Errorf("Discord Notify failed: %v", err)
	}
}
