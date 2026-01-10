package core

import (
	"strings"
	"time"
)

// Bounty represents a single unit of work
type Bounty struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Platform    string     `json:"platform"`
	Reward      string     `json:"reward"`
	Currency    string     `json:"currency"`
	URL         string     `json:"url"`
	CreatedAt   time.Time  `json:"created_at"`
	Score       int        `json:"score"`
	Description string     `json:"description"`
	Tags        []string   `json:"tags"`
	ExpiresAt   *time.Time `json:"expires_at"`
	PaymentType string     `json:"payment_type"`
}

// PaymentPriority defines the priority hierarchy
type PaymentPriority int

const (
	CryptoKing   PaymentPriority = iota // Highest priority
	P2PPremium                          // Cash App, Venmo
	FiatStandard                        // PayPal, Stripe, Wise
	LowPriority                         // Everything else
)

type PaymentConfig struct {
	CryptoCurrencies []string
	P2PMethods       []string
	FiatMethods      []string
}

var paymentConfig = defaultPaymentConfig()

func SetPaymentConfig(cfg PaymentConfig) {
	paymentConfig = normalizePaymentConfig(cfg)
}

// GetPaymentPriority returns the priority level for a given currency/payment type
func (b *Bounty) GetPaymentPriority() PaymentPriority {
	currency := b.Currency
	paymentType := b.PaymentType

	if containsCurrency(currency, paymentConfig.CryptoCurrencies) {
		return CryptoKing
	}
	if paymentType == "p2p" && containsCurrency(currency, paymentConfig.P2PMethods) {
		return P2PPremium
	}
	if containsCurrency(currency, paymentConfig.FiatMethods) {
		return FiatStandard
	}

	return LowPriority
}

func defaultPaymentConfig() PaymentConfig {
	return PaymentConfig{
		CryptoCurrencies: []string{"USDC", "USDT", "SOL", "ETH", "BTC", "MATIC", "AVAX", "ARB", "OP"},
		P2PMethods:       []string{"CASHAPP", "VENMO", "CASH APP"},
		FiatMethods:      []string{"USD", "PAYPAL", "STRIPE", "WISE"},
	}
}

func normalizePaymentConfig(cfg PaymentConfig) PaymentConfig {
	defaults := defaultPaymentConfig()
	cfg.CryptoCurrencies = normalizeUpperList(coalesceList(cfg.CryptoCurrencies, defaults.CryptoCurrencies))
	cfg.P2PMethods = normalizeUpperList(coalesceList(cfg.P2PMethods, defaults.P2PMethods))
	cfg.FiatMethods = normalizeUpperList(coalesceList(cfg.FiatMethods, defaults.FiatMethods))
	return cfg
}

func coalesceList(list []string, fallback []string) []string {
	if len(list) == 0 {
		return append([]string(nil), fallback...)
	}
	return list
}

func normalizeUpperList(list []string) []string {
	out := make([]string, 0, len(list))
	seen := make(map[string]struct{}, len(list))
	for _, item := range list {
		trimmed := strings.ToUpper(strings.TrimSpace(item))
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func containsCurrency(value string, tokens []string) bool {
	upper := strings.ToUpper(value)
	compact := strings.ReplaceAll(upper, " ", "")
	for _, token := range tokens {
		if token == "" {
			continue
		}
		target := strings.ToUpper(token)
		targetCompact := strings.ReplaceAll(target, " ", "")
		if strings.Contains(upper, target) || strings.Contains(compact, targetCompact) {
			return true
		}
	}
	return false
}
