CREATE TABLE IF NOT EXISTS saved (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_saved_user_id ON saved (user_id);
CREATE INDEX idx_saved_listing_id ON saved (listing_id);
CREATE UNIQUE INDEX idx_saved_user_listing_unique
    ON saved (user_id, listing_id)
    WHERE deleted_at IS NULL;

---- create above / drop below ----

DROP TABLE IF EXISTS saved;
