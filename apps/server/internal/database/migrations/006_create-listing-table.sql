CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand_id UUID NOT NULL REFERENCES brands(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    image_url TEXT[] NOT NULL,
    video_url TEXT[],
    condition listing_condition NOT NULL DEFAULT 'new',
    is_promoted BOOLEAN NOT NULL DEFAULT FALSE,
    is_discounted BOOLEAN NOT NULL DEFAULT FALSE,
    discount_percentage INTEGER NOT NULL DEFAULT 0 CHECK (discount_percentage >= 0 AND discount_percentage <= 100),
    brand_name VARCHAR(100),
    item_model VARCHAR(100),
    size VARCHAR(50),
    storage_size VARCHAR(50),
    color VARCHAR(50),
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_listings_brand_id ON listings (brand_id);
CREATE INDEX idx_listings_category_id ON listings (category_id);
CREATE INDEX idx_listings_condition ON listings (condition);
CREATE INDEX idx_listings_price ON listings (price);
CREATE INDEX idx_listings_is_promoted ON listings (is_promoted);
CREATE INDEX idx_listings_created_at ON listings (created_at DESC);

---- create above / drop below ----

DROP TABLE IF EXISTS listings;
