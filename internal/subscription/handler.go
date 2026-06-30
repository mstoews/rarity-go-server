package subscription

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstoews/rarity-go-server/internal/middleware"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /subscription/status
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFrom(r.Context())
	var status string
	var exp *time.Time
	h.db.QueryRow(r.Context(), `SELECT sub_status, sub_expires_at FROM users WHERE id=$1`, userID).
		Scan(&status, &exp)
	isActive := status == "active" && (exp == nil || exp.After(time.Now()))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"subscription_status": status,
		"is_active":           isActive,
		"sub_expires_at":      exp,
	})
}

// POST /subscription/verify  — receives StoreKit JWS transaction
func (h *Handler) Verify(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserIDFrom(r.Context())
	var body struct {
		SignedTransaction string `json:"signed_transaction"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	// TODO: verify JWS with Apple's App Store cert chain
	// For now: trust the payload and mark subscription active for 1 year
	exp := time.Now().AddDate(1, 0, 0)
	h.db.Exec(r.Context(),
		`UPDATE users SET sub_status='active', sub_expires_at=$1 WHERE id=$2`, exp, userID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"subscription_status": "active",
		"is_active":           true,
		"sub_expires_at":      exp,
	})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
