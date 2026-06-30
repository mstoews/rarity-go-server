package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstoews/rarity-go-server/internal/auth"
	"github.com/mstoews/rarity-go-server/internal/categories"
	"github.com/mstoews/rarity-go-server/internal/cosmetics"
	"github.com/mstoews/rarity-go-server/internal/middleware"
	"github.com/mstoews/rarity-go-server/internal/reviews"
	"github.com/mstoews/rarity-go-server/internal/stores"
	"github.com/mstoews/rarity-go-server/internal/subscription"
	"github.com/mstoews/rarity-go-server/internal/upload"
	"github.com/mstoews/rarity-go-server/internal/wishlist"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://rarity:rarity@localhost:5433/rarity?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	authH := auth.NewHandler(pool)
	catH := categories.NewHandler(pool)
	cosH := cosmetics.NewHandler(pool)
	storeH := stores.NewHandler(pool)
	revH := reviews.NewHandler(pool)
	wlH := wishlist.NewHandler(pool)
	subH := subscription.NewHandler(pool)
	uploadH, err := upload.NewHandler()
	if err != nil {
		log.Fatalf("upload handler: %v", err)
	}

	mux := http.NewServeMux()

	// Auth
	mux.HandleFunc("POST /auth/register", authH.Register)
	mux.HandleFunc("POST /auth/login",    authH.Login)
	mux.HandleFunc("POST /auth/apple",    authH.Apple)
	mux.HandleFunc("POST /auth/refresh",  authH.Refresh)

	// Protected helpers
	secret := authH.Secret()
	guard := func(h http.HandlerFunc) http.Handler {
		return middleware.RequireAuth(secret, h)
	}

	// Categories (public)
	mux.HandleFunc("GET /categories", catH.List)

	// Cosmetics
	mux.Handle("GET /cosmetics",      guard(cosH.List))
	mux.Handle("GET /cosmetics/{id}", guard(cosH.Detail))

	// Stores
	mux.Handle("GET /stores",      guard(storeH.List))
	mux.Handle("GET /stores/{id}", guard(storeH.Detail))

	// Reviews
	mux.Handle("GET  /cosmetics/{id}/reviews",  guard(revH.List))
	mux.Handle("POST /cosmetics/{id}/reviews",  guard(revH.Add))
	mux.Handle("DELETE /reviews/{id}",          guard(revH.Delete))

	// Wishlist
	mux.Handle("GET    /wishlist",     guard(wlH.List))
	mux.Handle("POST   /wishlist/{id}", guard(wlH.Add))
	mux.Handle("DELETE /wishlist/{id}", guard(wlH.Remove))

	// Subscription
	mux.Handle("GET  /subscription/status", guard(subH.Status))
	mux.Handle("POST /subscription/verify", guard(subH.Verify))

	// Upload
	mux.Handle("POST /upload/presigned-url", guard(uploadH.PresignedURL))

	addr := ":8092"
	log.Printf("rarity-go-server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, cors(mux)))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
