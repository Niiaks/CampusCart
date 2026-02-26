package model

// Category defines all fields for the category table in db.
type Category struct {
	ID       string `json:"category_id"`
	Name     string `json:"name" validate:"required"`
	ImageUrl string `json:"image_url" validate:"required"`
	PublicID string `json:"public_id" validate:"required"`
	Model
}
