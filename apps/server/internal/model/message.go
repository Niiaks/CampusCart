package model

import "time"

// Message represents a single message exchanged between users in a conversation.
type Message struct {
	ID             string     `json:"id"`                                                          // unique message id
	ConversationID string     `json:"conversation_id" validate:"required"`                         // thread this message belongs to
	SenderID       string     `json:"sender_id" validate:"required"`                               // author user id
	RecipientID    string     `json:"recipient_id" validate:"required"`                            // target user id
	ListingID      string     `json:"listing_id,omitempty"`                                        // optional listing context
	Content        string     `json:"content" validate:"required_without=MediaUrl"`                // text body
	MediaUrl       string     `json:"media_url,omitempty"`                                         // attached media link
	MediaType      string     `json:"media_type,omitempty" validate:"omitempty,oneof=photo video"` // media kind
	IsRead         bool       `json:"is_read"`                                                     // read status flag
	ReadAt         *time.Time `json:"read_at,omitempty"`                                           // when recipient read
	SentAt         time.Time  `json:"sent_at"`                                                     // when sent
	Model
}

// Conversation represents a message thread between a buyer and a seller, optionally tied to a listing.
type Conversation struct {
	ID            string    `json:"id"`                     // unique conversation id
	ListingID     string    `json:"listing_id,omitempty"`   // optional listing link
	BuyerID       string    `json:"buyer_id"`               // participant buyer id
	SellerID      string    `json:"seller_id"`              // participant seller id
	LastMessage   string    `json:"last_message,omitempty"` // cached last message text
	LastMessageAt time.Time `json:"last_message_at"`        // timestamp of last activity
	Model
}
