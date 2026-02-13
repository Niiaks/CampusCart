CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    type feedback_type NOT NULL,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedback_user_id ON feedback (user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_feedback_type ON feedback (type) WHERE deleted_at IS NULL;

---- create above / drop below ----

DROP TABLE IF EXISTS feedback;
