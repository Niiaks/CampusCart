CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id UUID REFERENCES listings(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    media_url TEXT,
    media_type media_type,
    CHECK ((media_url IS NULL AND media_type IS NULL) OR (media_url IS NOT NULL AND media_type IS NOT NULL)),
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    read_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_messages_sender_id ON messages (sender_id);
CREATE INDEX idx_messages_recipient_id ON messages (recipient_id);
CREATE INDEX idx_messages_listing_id ON messages (listing_id);
CREATE INDEX idx_messages_sent_at ON messages (sent_at DESC);
CREATE INDEX idx_messages_conversation ON messages (sender_id, recipient_id, listing_id);

---- create above / drop below ----

DROP TABLE IF EXISTS messages;
