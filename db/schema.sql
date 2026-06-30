-- Rarity PostgreSQL schema

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT UNIQUE,
    apple_sub       TEXT UNIQUE,
    username        TEXT NOT NULL UNIQUE,
    password_hash   TEXT,
    avatar_url      TEXT,
    sub_status      TEXT NOT NULL DEFAULT 'free'   CHECK (sub_status IN ('free','active','expired','cancelled')),
    sub_expires_at  TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Categories (admin-managed)
CREATE TABLE categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Cosmetics (admin-curated)
CREATE TABLE cosmetics (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    brand           TEXT NOT NULL,
    tagline         TEXT,
    description     TEXT,
    ingredients     TEXT,
    image_url       TEXT,
    images          TEXT[]   NOT NULL DEFAULT '{}',
    category_id     UUID REFERENCES categories(id),
    avg_rating      NUMERIC(3,2) NOT NULL DEFAULT 0,
    review_count    INT NOT NULL DEFAULT 0,
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Stores
CREATE TABLE stores (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    address         TEXT,
    city            TEXT,
    latitude        DOUBLE PRECISION,
    longitude       DOUBLE PRECISION,
    website         TEXT,
    opening_hours   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Cosmetic ↔ Store junction
CREATE TABLE cosmetic_stores (
    cosmetic_id UUID NOT NULL REFERENCES cosmetics(id) ON DELETE CASCADE,
    store_id    UUID NOT NULL REFERENCES stores(id)    ON DELETE CASCADE,
    in_stock    BOOLEAN NOT NULL DEFAULT TRUE,
    notes       TEXT,
    PRIMARY KEY (cosmetic_id, store_id)
);

-- Reviews
CREATE TABLE reviews (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cosmetic_id UUID NOT NULL REFERENCES cosmetics(id) ON DELETE CASCADE,
    user_id     UUID NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
    rating      SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    text        TEXT,
    photo_url   TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cosmetic_id, user_id)
);

-- Wishlist
CREATE TABLE wishlist (
    user_id     UUID NOT NULL REFERENCES users(id)      ON DELETE CASCADE,
    cosmetic_id UUID NOT NULL REFERENCES cosmetics(id)  ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, cosmetic_id)
);

-- Triggers: keep avg_rating / review_count in sync
CREATE OR REPLACE FUNCTION update_cosmetic_rating() RETURNS TRIGGER LANGUAGE plpgsql AS $$
BEGIN
    UPDATE cosmetics
    SET review_count = (SELECT COUNT(*) FROM reviews WHERE cosmetic_id = COALESCE(NEW.cosmetic_id, OLD.cosmetic_id)),
        avg_rating   = COALESCE((SELECT AVG(rating) FROM reviews WHERE cosmetic_id = COALESCE(NEW.cosmetic_id, OLD.cosmetic_id)), 0)
    WHERE id = COALESCE(NEW.cosmetic_id, OLD.cosmetic_id);
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_review_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW EXECUTE FUNCTION update_cosmetic_rating();
