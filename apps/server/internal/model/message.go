package model

import "time"

// Message defines all fields for the message table in db.
// for conversations
type Message struct {
	ID          string    `json:"id"`
	SenderID    string    `json:"sender_id" validate:"required"`
	RecipientID string    `json:"recipient_id" validate:"required"`
	ListingID   string    `json:"listing_id,omitempty"`
	Content     string    `json:"content" validate:"required"`
	MediaUrl    string    `json:"media_url,omitempty"`
	MediaType   string    `json:"media_type" validate:"required,oneof=photo video"`
	SentAt      time.Time `json:"sent_at"`
	ReadAt      time.Time `json:"read_at"`
	Model
}
