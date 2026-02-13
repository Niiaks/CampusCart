package model

// Feedback defines all fields for the feedback table in db.
type Feedback struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id" validate:"required"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=suggestion bug"`
	Model
}
