package reviews

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstoews/rarity-go-server/internal/middleware"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /cosmetics/{id}/reviews
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	cosmeticID := r.PathValue("id")
	cursor := r.URL.Query().Get("cursor")

	sql := `SELECT r.id, r.rating, r.text, r.photo_url, r.created_at, u.id, u.username
	        FROM reviews r JOIN users u ON u.id = r.user_id
	        WHERE r.cosmetic_id = $1`
	args := []any{cosmeticID}
	if cursor != "" {
		sql += ` AND r.id > $2`
		args = append(args, cursor)
	}
	sql += ` ORDER BY r.created_at DESC LIMIT 21`

	rows, err := h.db.Query(r.Context(), sql, args...)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var revs []map[string]any
	for rows.Next() {
		var rid, rating, uid, uname string
		var text, photoURL *string
		var createdAt string
		rows.Scan(&rid, &rating, &text, &photoURL, &createdAt, &uid, &uname)
		revs = append(revs, map[string]any{
			"id": rid, "rating": rating, "text": text, "photo_url": photoURL,
			"created_at": createdAt,
			"user":       map[string]any{"id": uid, "username": uname},
		})
	}

	var nextCursor *string
	if len(revs) == 21 {
		last := revs[20]["id"].(string)
		nextCursor = &last
		revs = revs[:20]
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"reviews": revs, "next_cursor": nextCursor})
}

// POST /cosmetics/{id}/reviews
func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {
	cosmeticID := r.PathValue("id")
	userID := middleware.UserIDFrom(r.Context())
	if userID == "" {
		jsonErr(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		Rating   int     `json:"rating"`
		Text     *string `json:"text"`
		PhotoURL *string `json:"photo_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	var rid, createdAt string
	err := h.db.QueryRow(r.Context(), `
		INSERT INTO reviews (cosmetic_id, user_id, rating, text, photo_url)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (cosmetic_id, user_id) DO UPDATE SET rating=$3, text=$4, photo_url=$5
		RETURNING id, created_at`,
		cosmeticID, userID, body.Rating, body.Text, body.PhotoURL,
	).Scan(&rid, &createdAt)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	var uname string
	h.db.QueryRow(r.Context(), `SELECT username FROM users WHERE id=$1`, userID).Scan(&uname)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id": rid, "rating": body.Rating, "text": body.Text, "photo_url": body.PhotoURL,
		"created_at": createdAt,
		"user":       map[string]any{"id": userID, "username": uname},
	})
}

// DELETE /reviews/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	rid := r.PathValue("id")
	userID := middleware.UserIDFrom(r.Context())
	h.db.Exec(r.Context(), `DELETE FROM reviews WHERE id=$1 AND user_id=$2`, rid, userID)
	w.WriteHeader(http.StatusNoContent)
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
