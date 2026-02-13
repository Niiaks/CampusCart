package model

// Brand defines all fields for the brand table in db.
type Brand struct {
	ID          string `json:"id"`
	SellerID    string `json:"seller_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	ProfileUrl  string `json:"profile_url,omitempty"`
	BannerUrl   string `json:"banner_url,omitempty"`
	Model
}
