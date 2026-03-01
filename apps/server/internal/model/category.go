// Package model defines data structures for representing categories and their attributes within the database.
//
// The Category struct models hierarchical item organization, supporting nested categories via ParentID,
// unique slugs for URL mapping, icons, sort order, and activation status. It also tracks creation and update timestamps.
// Children and Attributes fields are populated on demand to represent subcategories and category-specific metadata.
//
// The CategoryAttribute struct describes metadata fields associated with categories, including their type,
// display label, options (for enumerated types), requirement status, and sort order. OptionsRaw stores the raw
// database value for options, while Options is populated from JSONB for use in application logic.
// Category represents a category entity in the database, used for organizing items hierarchically.
// Each category may have a parent (for nested categories), a unique slug for URL mapping, an icon,
// and a sort order for display purposes. The IsActive field indicates whether the category is currently enabled.
// CreatedAt and UpdatedAt track the timestamps for creation and last modification.
// The Children and Attributes fields are populated on demand and not stored directly in the database.
// Children holds nested subcategories, while Attributes contains additional metadata specific to the category.
package model

import (
	"encoding/json"
	"time"
)

// Attribute type constants
const (
	AttributeTypeText    = "text"
	AttributeTypeNumber  = "number"
	AttributeTypeBoolean = "boolean"
	AttributeTypeEnum    = "enum"
)

type Category struct {
	ID        string  `json:"id"`                  // unique category id
	ParentID  *string `json:"parent_id,omitempty"` // nil = top-level category
	Name      string  `json:"name"`                // display label
	Slug      string  `json:"slug"`                // url-safe identifier
	Icon      string  `json:"icon,omitempty"`      // optional icon name/url
	PublicID  string  `json:"public_id"`           // id for cloudinary deletion
	IsActive  bool    `json:"is_active"`           // controls visibility
	SortOrder int     `json:"sort_order"`          // ordering among siblings
	Model

	// Populated on demand (not from db directly)
	Children   []Category          `json:"children,omitempty"`
	Attributes []CategoryAttribute `json:"attributes,omitempty"`
}

type CategoryAttribute struct {
	ID         string          `json:"id"`                // unique attribute id
	CategoryID string          `json:"category_id"`       // owning category
	Name       string          `json:"name"`              // key used in attributes map
	Label      string          `json:"label"`             // human-friendly label
	Type       string          `json:"type"`              // text|number|boolean|enum
	Options    []string        `json:"options,omitempty"` // allowed values for enum
	OptionsRaw json.RawMessage `json:"-"`                 // raw db value
	Required   bool            `json:"required"`          // whether seller must fill it
	SortOrder  int             `json:"sort_order"`        // ordering on forms
	CreatedAt  time.Time       `json:"created_at"`        // audit timestamp
}
