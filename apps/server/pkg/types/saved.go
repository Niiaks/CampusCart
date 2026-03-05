package types

type Save struct {
	UserID    string `json:"user_id" validate:"required"`
	ListingID string `json:"listing_id" validate:"required"`
}

func (s *Save) Validate() error { return nil }
