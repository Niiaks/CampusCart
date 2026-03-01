// Model defines common timestamp fields for database entities.
// It includes CreatedAt and UpdatedAt to track when records are created and modified,
// and DeletedAt as a nullable field to support soft deletion functionality.
// All fields are serialized to JSON with snake_case keys.
package model

import "time"

type Model struct {
	DeletedAt *time.Time `json:"deleted_at"` // null when active; set for soft deletes
	CreatedAt time.Time  `json:"created_at"` // record creation time
	UpdatedAt time.Time  `json:"updated_at"` // last update time
}
