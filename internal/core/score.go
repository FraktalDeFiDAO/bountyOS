package core

import (
	"strings"
	"time"
)

type ScoringConfig struct {
	UrgencyKeywords    []string
	DevTaskKeywords    []string
	AutomationKeywords []string
	SecurityKeywords   []string
	AuditKeywords      []string
}

var scoringConfig = defaultScoringConfig()

func SetScoringConfig(cfg ScoringConfig) {
	scoringConfig = normalizeScoringConfig(cfg)
}

func defaultScoringConfig() ScoringConfig {
	return ScoringConfig{
		UrgencyKeywords:    []string{"URGENT", "ASAP", "CRITICAL", "IMMEDIATE", "EMERGENCY"},
		DevTaskKeywords:    []string{"FIX", "BUG", "API", "INTEGRATION", "SMART CONTRACT", "BLOCKCHAIN"},
		AutomationKeywords: []string{"SCRIPT", "BOT"},
		SecurityKeywords:   []string{"SECURITY", "VULNERABILITY", "PENTEST", "HACK", "EXPLOIT"},
		AuditKeywords:      []string{"AUDIT"},
	}
}

func normalizeScoringConfig(cfg ScoringConfig) ScoringConfig {
	defaults := defaultScoringConfig()
	cfg.UrgencyKeywords = normalizeUpperList(coalesceList(cfg.UrgencyKeywords, defaults.UrgencyKeywords))
	cfg.DevTaskKeywords = normalizeUpperList(coalesceList(cfg.DevTaskKeywords, defaults.DevTaskKeywords))
	cfg.AutomationKeywords = normalizeUpperList(coalesceList(cfg.AutomationKeywords, defaults.AutomationKeywords))
	cfg.SecurityKeywords = normalizeUpperList(coalesceList(cfg.SecurityKeywords, defaults.SecurityKeywords))
	cfg.AuditKeywords = normalizeUpperList(coalesceList(cfg.AuditKeywords, defaults.AuditKeywords))
	cfg.DevTaskKeywords = removeOverlap(cfg.DevTaskKeywords, cfg.AutomationKeywords)
	cfg.SecurityKeywords = removeOverlap(cfg.SecurityKeywords, cfg.AuditKeywords)
	return cfg
}

func removeOverlap(base []string, blacklist []string) []string {
	if len(base) == 0 || len(blacklist) == 0 {
		return base
	}
	deny := make(map[string]struct{}, len(blacklist))
	for _, item := range blacklist {
		deny[item] = struct{}{}
	}
	out := make([]string, 0, len(base))
	for _, item := range base {
		if _, ok := deny[item]; ok {
			continue
		}
		out = append(out, item)
	}
	return out
}

func containsAny(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if keyword != "" && strings.Contains(text, keyword) {
			return true
		}
	}
	return false
}

// CalculateUrgency applies the "Obsidian" scoring algorithm
func CalculateUrgency(b *Bounty) int {
	score := 0
	titleUpper := strings.ToUpper(b.Title)

	// ------------------------------------------
	// RULE 1: PAYMENT METHOD HIERARCHY (Crypto is King)
	// ------------------------------------------

	// TIER 0: KING CRYPTO (Instant Settlement)
	// We look for Stablecoins and Layer 1 tokens
	if containsCurrency(b.Currency, paymentConfig.CryptoCurrencies) {
		score += 50
	} else if containsCurrency(b.Currency, paymentConfig.P2PMethods) {
		// TIER 1: P2P FIAT (High Velocity)
		score += 45
	} else if containsCurrency(b.Currency, paymentConfig.FiatMethods) {
		// TIER 2: LEGACY FIAT (Medium Velocity)
		score += 25
	} else {
		// TIER 3: UNKNOWN / SLOW
		score += 5
	}

	// ------------------------------------------
	// RULE 2: KEYWORD TRIGGERS
	// ------------------------------------------
	if containsAny(titleUpper, scoringConfig.UrgencyKeywords) {
		score += 30
	}
	if containsAny(titleUpper, scoringConfig.DevTaskKeywords) {
		score += 15 // Dev tasks are usually quick
	}
	if containsAny(titleUpper, scoringConfig.AutomationKeywords) {
		score += 20 // Automation tasks (High value for you)
	}
	if containsAny(titleUpper, scoringConfig.SecurityKeywords) {
		score += 25 // Security tasks (High value)
	}
	if containsAny(titleUpper, scoringConfig.AuditKeywords) {
		score += 35 // Audit tasks (Very high value)
	}

	// ------------------------------------------
	// RULE 3: RECENCY (The Sniper Rule)
	// ------------------------------------------
	duration := time.Since(b.CreatedAt)
	if duration < 1*time.Hour {
		score += 40 // Super Fresh
	} else if duration < 6*time.Hour {
		score += 25 // Fresh
	} else if duration < 24*time.Hour {
		score += 10 // Recent
	}

	// ------------------------------------------
	// RULE 4: PLATFORM PRIORITY
	// ------------------------------------------
	platformUpper := strings.ToUpper(b.Platform)
	if strings.Contains(platformUpper, "SUPERTEAM") {
		score += 15 // Solana ecosystem, high value
	}
	if strings.Contains(platformUpper, "BOUNTYCASTER") {
		score += 10 // Social feed, fast payment
	}
	if strings.Contains(platformUpper, "IMMUNEFI") || strings.Contains(platformUpper, "HACKEN") {
		score += 30 // Bug bounty, high value
	}

	// Apply tags bonuses
	for _, tag := range b.Tags {
		tagUpper := strings.ToUpper(tag)
		if strings.Contains(tagUpper, "URGENT") {
			score += 20
		}
		if strings.Contains(tagUpper, "HOT") {
			score += 15
		}
		if strings.Contains(tagUpper, "DEADLINE") {
			score += 10
		}
	}

	return score
}
