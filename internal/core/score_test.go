package core

import (
	"testing"
	"time"
)

func TestCalculateUrgency(t *testing.T) {
	tests := []struct {
		name     string
		bounty   Bounty
		minScore int
		maxScore int
	}{
		{
			name: "Crypto King - USDC & Urgent",
			bounty: Bounty{
				Title:     "Urgent: Fix Security Bug",
				Currency:  "USDC",
				CreatedAt: time.Now(),
				Platform:  "GitHub",
			},
			minScore: 100, // 50 (Crypto) + 30 (Urgent) + 15 (Fix/Bug) + 25 (Security) + 40 (Fresh)
			maxScore: 200,
		},
		{
			name: "Fiat Standard - PayPal",
			bounty: Bounty{
				Title:     "Write a blog post",
				Currency:  "PAYPAL",
				CreatedAt: time.Now().Add(-2 * time.Hour), // Fresh (< 6h)
				Platform:  "GitHub",
			},
			minScore: 50, // 25 (Fiat) + 25 (Fresh)
			maxScore: 80,
		},
		{
			name: "P2P Premium - CashApp",
			bounty: Bounty{
				Title:     "Script needed",
				Currency:  "CASHAPP",
				CreatedAt: time.Now().Add(-20 * time.Hour), // Recent (< 24h)
				Platform:  "GitHub",
			},
			minScore: 65, // 45 (P2P) + 20 (Script) + 10 (Recent)
			maxScore: 90,
		},
		{
			name: "Superteam Bonus",
			bounty: Bounty{
				Title:     "Solana integration",
				Currency:  "SOL",
				CreatedAt: time.Now(),
				Platform:  "SUPERTEAM",
			},
			minScore: 105, // 50 (Crypto) + 40 (Fresh) + 15 (Superteam)
			maxScore: 150,
		},
		{
			name: "Audit Task High Value",
			bounty: Bounty{
				Title:     "Smart Contract Audit",
				Currency:  "USDT",
				CreatedAt: time.Now(),
				Platform:  "ImmuneFi",
			},
			minScore: 130, // 50 (Crypto) + 35 (Audit) + 40 (Fresh) + 30 (ImmuneFi)
			maxScore: 180,
		},
		{
			name: "Unknown Currency & Automation",
			bounty: Bounty{
				Title:     "Write a Python Script",
				Currency:  "ROCKS",
				CreatedAt: time.Now().Add(-5 * time.Hour), // Fresh < 6h
				Platform:  "Freelancer",
			},
			minScore: 50, // 5 (Unknown) + 20 (Script) + 25 (Fresh)
			maxScore: 70,
		},
		{
			name: "Tag Bonuses: Hot & Deadline",
			bounty: Bounty{
				Title:     "Simple Task",
				Currency:  "USDC",
				CreatedAt: time.Now(),
				Platform:  "GitHub",
				Tags:      []string{"hot", "deadline"},
			},
			minScore: 115, // 50 (Crypto) + 40 (Fresh) + 15 (Hot) + 10 (Deadline)
			maxScore: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateUrgency(&tt.bounty)
			if score < tt.minScore {
				t.Errorf("CalculateUrgency() score = %v, want >= %v", score, tt.minScore)
			}
		})
	}
}

func TestGetPaymentPriority(t *testing.T) {
	tests := []struct {
		name     string
		bounty   Bounty
		expected PaymentPriority
	}{
		{
			name: "USDC is CryptoKing",
			bounty: Bounty{Currency: "USDC", PaymentType: "crypto"},
			expected: CryptoKing,
		},
		{
			name: "CashApp is P2PPremium",
			bounty: Bounty{Currency: "CASHAPP", PaymentType: "p2p"},
			expected: P2PPremium,
		},
		{
			name: "PayPal is FiatStandard",
			bounty: Bounty{Currency: "PAYPAL", PaymentType: "fiat"},
			expected: FiatStandard,
		},
		{
			name: "Unknown is LowPriority",
			bounty: Bounty{Currency: "ROCKS", PaymentType: "barter"},
			expected: LowPriority,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := tt.bounty.GetPaymentPriority()
			if priority != tt.expected {
				t.Errorf("GetPaymentPriority() = %v, want %v", priority, tt.expected)
			}
		})
	}
}
