package types

type UpdateCategory struct {
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type CategoryResponse struct {
	ID       string `json:"category_id"`
	Name     string `json:"name" validate:"required"`
	ImageUrl string `json:"image_url" validate:"required"`
	PublicID string `json:"public_id" validate:"required"`
}
