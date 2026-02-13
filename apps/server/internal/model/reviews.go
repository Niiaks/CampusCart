package model

// Review defines all fields for the reviews table in db.
type Review struct {
	ID         string `json:"id"`
	ListingID  string `json:"listing_id" validate:"required"`
	ReviewerID string `json:"reviewer_id" validate:"required"`
	SellerID   string `json:"seller_id" validate:"required"`
	Rating     int    `json:"rating" validate:"required,min=1,max=5"`
	Comment    string `json:"comment,omitempty"`
	Model
}
