package model

// Listing defines all fields for the listing/ads table in db.
type Listing struct {
	ID                 string   `json:"id"`
	BrandID            string   `json:"brand_id" validate:"required"`
	CategoryID         string   `json:"category_id" validate:"required"`
	Title              string   `json:"title" validate:"required"`
	Description        string   `json:"description" validate:"required"`
	Price              int      `json:"price" validate:"required,gte=0"`
	ImageUrl           []string `json:"image_url" validate:"required"`
	VideoUrl           []string `json:"video_url,omitempty"`
	Condition          string   `json:"condition" validate:"required,oneof=new used second-hand"`
	IsPromoted         bool     `json:"is_promoted,omitempty"`
	IsDiscounted       bool     `json:"is_discounted,omitempty"`
	DiscountPercentage int      `json:"discount_percentage" validate:"gte=0"`
	Brand              string   `json:"brand,omitempty"`
	ItemModel          string   `json:"model,omitempty"`
	Size               string   `json:"size,omitempty"`
	StorageSize        string   `json:"storage_size,omitempty"`
	Color              string   `json:"color,omitempty"`
	Model
}
