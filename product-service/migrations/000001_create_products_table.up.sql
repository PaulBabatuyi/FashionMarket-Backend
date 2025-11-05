CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,  
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    image_url TEXT NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    category TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,

    tsv tsvector GENERATED ALWAYS AS (
        setweight(to_tsvector('english', name), 'A') ||
        setweight(to_tsvector('english', coalesce(description, '')), 'B')
    ) STORED
);


CREATE INDEX idx_products_user_id ON products(user_id);
CREATE INDEX idx_products_category ON products USING GIN(category);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_products_tsv ON products USING GIN(tsv);
--CREATE INDEX idx_products_category_created ON products USING GIN(category) 
  --  WHERE created_at > NOW() - INTERVAL '30 days';  