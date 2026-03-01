package model

// Listing represents an advertisement or item listing within the application.
// It contains all relevant fields required for storing and retrieving listing data from the database.
import "time"

type Listing struct {
	ID          string                 `json:"id"`                                                       // unique listing id
	SellerID    string                 `json:"seller_id"`                                                // owner user id
	CategoryID  string                 `json:"category_id" validate:"required"`                          // category placement
	Title       string                 `json:"title" validate:"required,min=5,max=255"`                  // headline shown to buyers
	Description string                 `json:"description" validate:"required"`                          // full description
	Price       int64                  `json:"price" validate:"required,gte=0"`                          // amount in minor units
	Condition   string                 `json:"condition" validate:"required,oneof=new used second-hand"` // item condition
	Negotiable  bool                   `json:"negotiable"`                                               // true if price is flexible
	Attributes  map[string]interface{} `json:"attributes,omitempty"`                                     // dynamic category-specific fields
	ImageUrls   []string               `json:"image_urls" validate:"required,min=1"`                     // photos for the listing
	VideoUrls   []string               `json:"video_urls,omitempty"`                                     // optional video links
	IsActive    bool                   `json:"is_active"`                                                // controls visibility
	IsPromoted  bool                   `json:"is_promoted"`                                              // boosted placement flag
	ViewsCount  int                    `json:"views_count"`                                              // view counter
	CreatedAt   time.Time              `json:"created_at"`                                               // created timestamp
	UpdatedAt   time.Time              `json:"updated_at"`                                               // last updated timestamp

	// Populated on demand (joins)
	Category *Category `json:"category,omitempty"` // populated via joins
	Seller   *User     `json:"seller,omitempty"`   // populated via joins
}
