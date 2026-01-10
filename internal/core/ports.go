package core

import "context"

// Scanner interface for data source integrations
type Scanner interface {
	Name() string
	Scan(ctx context.Context) (<-chan Bounty, error)
}

// Notifier interface for alerting systems
type Notifier interface {
	Alert(bounty Bounty) error
	Notify(message string) error
}

// Storage interface for persistence
type Storage interface {
	Save(bounty Bounty) error
	IsNew(url string) (bool, error)
	GetRecent(limit int) ([]Bounty, error)
	Close() error
}