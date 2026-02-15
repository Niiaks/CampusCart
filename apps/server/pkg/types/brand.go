package types

type UpdateBrand struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	ProfileUrl  *string `json:"profile_url"`
	BannerUrl   *string `json:"banner_url"`
}
