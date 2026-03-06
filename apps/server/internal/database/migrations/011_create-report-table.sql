-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type report_entity_type NOT NULL,
    entity_id UUID NOT NULL,
    reason report_reason NOT NULL,
    description TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    admin_note TEXT,
    reviewed_by UUID REFERENCES users(id),
    reviewed_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_reporter_id ON reports (reporter_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reports_entity ON reports (entity_type, entity_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_reports_status ON reports (status) WHERE deleted_at IS NULL;
CREATE INDEX idx_reports_created_at ON reports (created_at DESC) WHERE deleted_at IS NULL;

-- Prevent duplicate reports from same user for same entity
CREATE UNIQUE INDEX idx_reports_unique_active ON reports (reporter_id, entity_type, entity_id) 
    WHERE deleted_at IS NULL AND status IN ('pending', 'under_review');

CREATE TRIGGER trg_reports_set_updated_at
    BEFORE UPDATE ON reports
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
    
---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_reports_set_updated_at ON reports;
DROP TABLE IF EXISTS reports;
DROP TYPE IF EXISTS report_reason;
DROP TYPE IF EXISTS report_entity_type;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.









---- create above / drop below ----

