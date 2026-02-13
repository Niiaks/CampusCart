package model

// Saved defines all fields for the saved table in db.

type Saved struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id" validate:"required"`
	ListingID string `json:"listing_id" validate:"required"`
	Model
}
