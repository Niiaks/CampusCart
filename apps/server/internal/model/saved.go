package model

// Saved represents a record in the "saved" table, linking a user to a saved listing.
// It contains the unique identifier for the saved entry, the user's ID, and the listing's ID.
// The struct embeds Model for common fields such as timestamps.
// Fields UserID and ListingID are required and validated accordingly.
type Saved struct {
	ID        string `json:"id"`                             // unique saved-row id
	UserID    string `json:"user_id" validate:"required"`    // owner user id
	ListingID string `json:"listing_id" validate:"required"` // saved listing id
	Model
}
