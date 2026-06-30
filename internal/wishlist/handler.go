package wishlist

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstoews/rarity-go-server/internal/middleware"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /wishlist
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFrom(r.Context())
	rows, err := h.db.Query(r.Context(), `
		SELECT c.id, c.name, c.brand, c.tagline, c.image_url, c.avg_rating, c.review_count
		FROM cosmetics c JOIN wishlist wl ON wl.cosmetic_id = c.id
		WHERE wl.user_id = $1 ORDER BY wl.created_at DESC`, userID)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var items []map[string]any
	for rows.Next() {
		var id, name, brand string
		var tagline, imageURL *string
		var avg float64
		var count int
		rows.Scan(&id, &name, &brand, &tagline, &imageURL, &avg, &count)
		items = append(items, map[string]any{
			"id": id, "name": name, "brand": brand, "tagline": tagline, "image_url": imageURL,
			"avg_rating": avg, "review_count": count,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"cosmetics": items})
}

// POST /wishlist/{id}
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	cosmeticID := r.PathValue("id")
	userID := middleware.UserIDFrom(r.Context())
	h.db.Exec(r.Context(), `INSERT INTO wishlist (user_id, cosmetic_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		userID, cosmeticID)
	w.WriteHeader(http.StatusNoContent)
}

// DELETE /wishlist/{id}
func (h *Handler) Remove(w http.ResponseWriter, r *http.Request) {
	cosmeticID := r.PathValue("id")
	userID := middleware.UserIDFrom(r.Context())
	h.db.Exec(r.Context(), `DELETE FROM wishlist WHERE user_id=$1 AND cosmetic_id=$2`, userID, cosmeticID)
	w.WriteHeader(http.StatusNoContent)
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
