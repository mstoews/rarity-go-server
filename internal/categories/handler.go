package categories

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /categories
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, name, description FROM categories WHERE TRUE ORDER BY sort_order, name`)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	cats := []map[string]any{}
	for rows.Next() {
		var id, name string
		var desc *string
		rows.Scan(&id, &name, &desc)
		cats = append(cats, map[string]any{"id": id, "name": name, "description": desc})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"categories": cats})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
