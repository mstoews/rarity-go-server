# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
# Start PostgreSQL via Docker (port 5433)
make postgres

# Apply schema
make migrate

# Run server (port 8092)
make server

# Or directly
go run ./cmd/server

# Build binary
go build ./...

# Run all tests
go test ./...
```

Environment variables (all have defaults for local dev):

| Variable | Default |
|---|---|
| `DATABASE_URL` | `postgres://rarity:rarity@localhost:5433/rarity?sslmode=disable` |
| `JWT_SECRET` | `dev-secret-change-in-prod` |
| `S3_BUCKET` | `rarity-uploads` |
| `AWS_*` | standard AWS SDK env vars for S3 uploads |

## Architecture

Single-binary Go REST API using the stdlib `net/http` mux (Go 1.22 `r.PathValue()`). No router framework. `pgx/v5` pgxpool for PostgreSQL. Port **8092**.

**Entry point:** `cmd/server/main.go` — wires all handlers onto the mux, wraps with CORS middleware.

**Package layout — one package per domain:**

| Package | Responsibility |
|---|---|
| `internal/auth` | Register, Login, Apple Sign-In, Refresh; JWT issue/parse; bcrypt passwords |
| `internal/middleware` | `RequireAuth` — Bearer token extraction + `UserIDFrom(ctx)` |
| `internal/categories` | Admin-curated category list (public endpoint) |
| `internal/cosmetics` | Cosmetic listing (free/paid gate) and full detail (paid only → 402) |
| `internal/stores` | Store list with lat/lng query and store detail with cosmetics |
| `internal/reviews` | Per-cosmetic reviews; upsert (one review per user per cosmetic) |
| `internal/wishlist` | User wishlist CRUD |
| `internal/subscription` | StoreKit 2 transaction verify + subscription status |
| `internal/upload` | S3 presigned PUT URL generation |

## Data Model

See `db/schema.sql` for the full schema. Key relationships:

- `cosmetic_stores` — many-to-many join between `cosmetics` and `stores`, carries `in_stock` and `notes`
- `reviews` — unique on `(cosmetic_id, user_id)`; a trigger on this table keeps `cosmetics.avg_rating` and `cosmetics.review_count` in sync
- `users.sub_status` — `'free' | 'active' | 'expired' | 'cancelled'`; `sub_expires_at` is the wall-clock expiry checked at request time

## Subscription Gate

The paywall is enforced in `internal/cosmetics/handler.go`. `GET /cosmetics` returns all cosmetics but omits `avg_rating` and `review_count` for non-subscribers. `GET /cosmetics/{id}` returns **402 Payment Required** for non-subscribers — the iOS client intercepts this and shows the paywall sheet.

Subscription status is checked inline via a DB query on every gated request (no caching). The `isSubscribed` helper queries `users` for `sub_status='active'` AND `sub_expires_at > NOW()`.

## Auth Flow

1. Client sends `POST /auth/login` or `/auth/register` → receives `access_token` (15 min HS256 JWT) + `refresh_token` (raw hex)
2. Refresh token is stored hashed (bcrypt) in `refresh_tokens` table
3. On 401, client calls `POST /auth/refresh` — server iterates non-expired tokens and bcrypt-compares to find a match, then rotates (delete old, issue new)
4. Apple Sign-In: `POST /auth/apple` — `verifyAppleToken` in `auth/handler.go` is a stub; replace with JWS verification against Apple's public keys before shipping

## Stubs to Complete Before Production

- `verifyAppleToken` in `internal/auth/handler.go` — needs real Apple JWS verification
- `internal/subscription/handler.go` `Verify` — needs JWS verification against Apple's App Store certificate chain
- `internal/stores/handler.go` `List` — lat/lng params are parsed but not yet used for distance ordering; add PostGIS or a Haversine ORDER BY
