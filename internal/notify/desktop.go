package notify

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"bountyos-v8/internal/core"
)

type DesktopNotifier struct{}

func NewDesktopNotifier() *DesktopNotifier {
	return &DesktopNotifier{}
}

func (n *DesktopNotifier) Alert(bounty core.Bounty) error {
	message := fmt.Sprintf("New Bounty: %s\nPlatform: %s\nReward: %s\nLink: %s", bounty.Title, bounty.Platform, bounty.Reward, bounty.URL)
	return n.notify(message, bounty.URL)
}

func (n *DesktopNotifier) Notify(message string) error {
	return n.notify(message, "")
}

func (n *DesktopNotifier) notify(message string, link string) error {
	// Check for headless mode (e.g., Docker)
	if os.Getenv("HEADLESS") == "true" {
		log.Printf("[NOTIFY] %s", message)
		return nil
	}

	var err error
	switch runtime.GOOS {
	case "linux":
		if link != "" {
			go n.notifyLinuxWithAction(message, link)
			return nil
		}
		err = exec.Command("notify-send", "BountyOS Alert", message).Run()
	case "darwin":
		if link != "" {
			if path, lookErr := exec.LookPath("terminal-notifier"); lookErr == nil {
				err = exec.Command(path, "-title", "BountyOS Alert", "-message", message, "-open", link).Run()
				break
			}
		}
		err = exec.Command("osascript", "-e", fmt.Sprintf(`display notification "%s" with title "BountyOS Alert"`, message)).Run()
	case "windows":
		// Windows notification would require additional libraries
		// For now, just print to console
		fmt.Println("BountyOS Alert:", message)
		return nil
	default:
		fmt.Println("BountyOS Alert:", message)
		return nil
	}

	if err != nil {
		// Fallback to log if notification fails (common in headless/container envs)
		log.Printf("[NOTIFY FAIL] %s (Error: %v)", message, err)
		return nil
	}
	return nil
}

func (n *DesktopNotifier) notifyLinuxWithAction(message string, link string) {
	cmd := exec.Command("notify-send", "--action=default=Open", "--wait", "BountyOS Alert", message)
	output, err := cmd.Output()
	if err != nil {
		_ = exec.Command("notify-send", "BountyOS Alert", message).Run()
		return
	}

	if strings.TrimSpace(string(output)) != "" {
		_ = openURL(link)
	}
}

func openURL(link string) error {
	if link == "" {
		return nil
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", link).Run()
	case "darwin":
		return exec.Command("open", link).Run()
	case "windows":
		return exec.Command("cmd", "/c", "start", "", link).Run()
	default:
		return nil
	}
}
