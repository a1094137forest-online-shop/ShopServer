CREATE TABLE shops (
    id BIGSERIAL PRIMARY KEY,
    shop_id VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(100) NOT NULL DEFAULT 'My Shop',
    account VARCHAR(100) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    product_id VARCHAR(30) UNIQUE NOT NULL,
    shop_id BIGINT NOT NULL REFERENCES shops (id),
    title VARCHAR(500) NOT NULL,
    description VARCHAR(500) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);