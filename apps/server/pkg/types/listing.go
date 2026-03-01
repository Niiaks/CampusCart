package types

type UpdateListing struct {
	Title       *string                 `json:"title"`
	Description *string                 `json:"description"`
	CategoryID  *string                 `json:"category_id"`
	Price       *int64                  `json:"price"`
	Condition   *string                 `json:"condition"`
	Negotiable  *bool                   `json:"negotiable"`
	Attributes  *map[string]interface{} `json:"attributes"`
	ImageUrls   *[]string               `json:"image_urls"`
	VideoUrls   *[]string               `json:"video_urls"`
	IsActive    *bool                   `json:"is_active"`
	IsPromoted  *bool                   `json:"is_promoted"`
}

type ListingFilter struct {
	CategoryID string
	BrandID    string
	Search     string
	MinPrice   *int64
	MaxPrice   *int64
	Condition  string
	Limit      int
	Offset     int
}
