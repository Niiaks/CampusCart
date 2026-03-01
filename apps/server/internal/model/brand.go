package model

// Brand represents a seller's brand entity in the database.
// It encapsulates all relevant information and identity fields for a seller
type Brand struct {
	ID          string `json:"id"`                                      // unique brand identifier
	SellerID    string `json:"seller_id" validate:"required"`           // owner user id
	Name        string `json:"name" validate:"required,min=2,max=100"`  // display name for the brand
	Slug        string `json:"slug"`                                    // url-friendly name, e.g., /shop/johns-store
	Description string `json:"description" validate:"required,max=500"` // short bio/about
	ProfileUrl  string `json:"profile_url,omitempty"`                   // avatar/logo image
	BannerUrl   string `json:"banner_url,omitempty"`                    // header/cover image
	IsVerified  bool   `json:"is_verified"`                             // shows verified seller badge
	Model
}
