package model

// Review represents a user's evaluation of a listing.
type Review struct {
	ID         string   `json:"id"`                                     // unique review id
	ListingID  string   `json:"listing_id" validate:"required"`         // listing being reviewed
	ReviewerID string   `json:"reviewer_id" validate:"required"`        // who wrote the review
	SellerID   string   `json:"seller_id" validate:"required"`          // seller being rated
	Rating     int      `json:"rating" validate:"required,min=1,max=5"` // star value 1-5
	Comment    string   `json:"comment,omitempty"`                      // optional text
	ImageUrls  []string `json:"image_urls,omitempty"`                   // optional photos
	Model
}
