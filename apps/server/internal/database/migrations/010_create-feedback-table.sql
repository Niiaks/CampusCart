CREATE TABLE IF NOT EXISTS feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    type feedback_type NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    admin_note TEXT,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedback_user_id ON feedback (user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_feedback_type ON feedback (type) WHERE deleted_at IS NULL;
CREATE INDEX idx_feedback_status ON feedback (status) WHERE deleted_at IS NULL;
CREATE TRIGGER trg_feedback_set_updated_at
    BEFORE UPDATE ON feedback
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_feedback_set_updated_at ON feedback;
DROP TABLE IF EXISTS feedback;
