package database

import "context"

// Pinger defines an interface for database health checking and lifecycle management
type Pinger interface {
	Ping(ctx context.Context) error
	Close() error
}

// Ensure Database implements Pinger
var _ Pinger = (*Database)(nil)

// Ping checks if the database is reachable
func (db *Database) Ping(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}
