package types

type UpdateCategory struct {
	Name      *string `json:"name"`
	ParentID  *string `json:"parent_id"`
	Slug      *string `json:"slug"`
	Icon      *string `json:"icon"`
	PublicID  *string `json:"public_id"`
	IsActive  *bool   `json:"is_active"`
	SortOrder *int    `json:"sort_order"`
}

type CategoryResponse struct {
	ID        string  `json:"category_id"`
	ParentID  *string `json:"parent_id,omitempty"`
	Name      string  `json:"name" validate:"required"`
	Slug      string  `json:"slug" validate:"required"`
	Icon      *string `json:"icon,omitempty"`
	PublicID  *string `json:"public_id,omitempty"`
	IsActive  bool    `json:"is_active"`
	SortOrder int     `json:"sort_order"`
}

type CategoryAttributeResponse struct {
	Name      string   `json:"name"`
	Label     string   `json:"label"`
	Type      string   `json:"type"`
	Options   []string `json:"options,omitempty"`
	Required  bool     `json:"required"`
	SortOrder int      `json:"sort_order"`
}
