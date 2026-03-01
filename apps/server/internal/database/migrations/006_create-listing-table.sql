CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    condition listing_condition NOT NULL DEFAULT 'new',
    negotiable BOOLEAN NOT NULL DEFAULT FALSE,
    attributes JSONB DEFAULT '{}'::jsonb,
    image_urls TEXT[] NOT NULL,
    video_urls TEXT[],
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_promoted BOOLEAN NOT NULL DEFAULT FALSE,
    views_count INTEGER NOT NULL DEFAULT 0,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE listings
    ADD CONSTRAINT listings_image_urls_not_empty CHECK (array_length(image_urls, 1) IS NOT NULL AND array_length(image_urls, 1) > 0);

CREATE INDEX idx_listings_brand_id ON listings (brand_id);
CREATE INDEX idx_listings_category_id ON listings (category_id);
CREATE INDEX idx_listings_condition ON listings (condition);
CREATE INDEX idx_listings_price ON listings (price);
CREATE INDEX idx_listings_is_promoted ON listings (is_promoted);
CREATE INDEX idx_listings_is_active ON listings (is_active);
CREATE INDEX idx_listings_created_at ON listings (created_at DESC);
CREATE TRIGGER trg_listings_set_updated_at
    BEFORE UPDATE ON listings
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_listings_set_updated_at ON listings;
DROP TABLE IF EXISTS listings;
