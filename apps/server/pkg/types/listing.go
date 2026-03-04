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
	IsPromoted  *bool                   `json:"-"`
}

// Validate implements validation.Validatable (lightweight; handler enforces business rules).
func (u *UpdateListing) Validate() error { return nil }

type CreateListingRequest struct {
	CategoryID  string                 `json:"category_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Price       int64                  `json:"price"`
	Condition   string                 `json:"condition"`
	Negotiable  bool                   `json:"negotiable"`
	Attributes  map[string]interface{} `json:"attributes"`
	ImageUrls   []string               `json:"image_urls"`
	VideoUrls   []string               `json:"video_urls"`
	IsActive    bool                   `json:"is_active"`
}

func (c *CreateListingRequest) Validate() error { return nil }

type ListingResponse struct {
	ID          string                 `json:"id"`
	BrandID     string                 `json:"brand_id"`
	CategoryID  string                 `json:"category_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Price       int64                  `json:"price"`
	Condition   string                 `json:"condition"`
	Negotiable  bool                   `json:"negotiable"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
	ImageUrls   []string               `json:"image_urls"`
	VideoUrls   []string               `json:"video_urls,omitempty"`
	IsActive    bool                   `json:"is_active"`
	IsPromoted  bool                   `json:"is_promoted"`
	ViewsCount  int                    `json:"views_count"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
}

type ListingFilter struct {
	CategoryID         string
	BrandID            string
	BrandName          string
	Search             string
	MinPrice           *int64
	MaxPrice           *int64
	Condition          string
	IncludeDescendants bool
	Limit              int
	Offset             int
}

func (f *ListingFilter) Validate() error { return nil }
