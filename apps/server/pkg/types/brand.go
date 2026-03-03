package types

type UpdateBrand struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ProfileUrl  *string `json:"profile_url"`
	BannerUrl   *string `json:"banner_url"`
}

// Validate implements validation.Validatable.
func (u *UpdateBrand) Validate() error { return nil }

type BrandResponse struct {
	ID          string  `json:"id"`
	SellerID    string  `json:"seller_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description,omitempty"`
	ProfileUrl  *string `json:"profile_url,omitempty"`
	BannerUrl   *string `json:"banner_url,omitempty"`
	IsVerified  bool    `json:"is_verified"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}
