-- Category hierarchy and attribute metadata
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(160) NOT NULL UNIQUE,
    icon TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent ON categories (parent_id);
CREATE INDEX idx_categories_sort ON categories (sort_order);
CREATE TRIGGER trg_categories_set_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

CREATE TABLE IF NOT EXISTS category_attributes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    label VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('text', 'number', 'boolean', 'enum')),
    options JSONB,
    required BOOLEAN NOT NULL DEFAULT FALSE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (category_id, name)
);

CREATE INDEX idx_category_attributes_category ON category_attributes (category_id);
CREATE INDEX idx_category_attributes_sort ON category_attributes (sort_order);
CREATE TRIGGER trg_category_attributes_set_updated_at
    BEFORE UPDATE ON category_attributes
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_category_attributes_set_updated_at ON category_attributes;
DROP TABLE IF EXISTS category_attributes;
DROP TRIGGER IF EXISTS trg_categories_set_updated_at ON categories;
DROP TABLE IF EXISTS categories;
