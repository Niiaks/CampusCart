package model

// Feedback represents a user's input to the system, which can be either a suggestion or a bug report.
type Feedback struct {
	ID        string `json:"id"`                                                  // unique feedback id
	UserID    string `json:"user_id" validate:"required"`                         // authoring user
	Message   string `json:"message" validate:"required,min=10,max=1000"`         // feedback body
	Type      string `json:"type" validate:"required,oneof=suggestion bug other"` // category of feedback
	Status    string `json:"status"`                                              // open | reviewed | resolved
	AdminNote string `json:"admin_note,omitempty"`                                // internal note from team
	Model
}
