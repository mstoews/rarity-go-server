package stores

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /stores?lat=&lng=
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat, _ := strconv.ParseFloat(q.Get("lat"), 64)
	lng, _ := strconv.ParseFloat(q.Get("lng"), 64)
	_ = lat
	_ = lng

	rows, err := h.db.Query(r.Context(),
		`SELECT id, name, address, city, latitude, longitude FROM stores ORDER BY name LIMIT 50`)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var stores []map[string]any
	for rows.Next() {
		var id, name string
		var addr, city *string
		var lat2, lng2 *float64
		rows.Scan(&id, &name, &addr, &city, &lat2, &lng2)
		stores = append(stores, map[string]any{
			"id": id, "name": name, "address": addr, "city": city,
			"latitude": lat2, "longitude": lng2,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"stores": stores})
}

// GET /stores/{id}
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var sid, name string
	var addr, city, website, hours *string
	var lat, lng *float64
	err := h.db.QueryRow(r.Context(),
		`SELECT id, name, address, city, latitude, longitude, website, opening_hours FROM stores WHERE id=$1`, id,
	).Scan(&sid, &name, &addr, &city, &lat, &lng, &website, &hours)
	if err != nil {
		jsonErr(w, "not found", http.StatusNotFound)
		return
	}

	cosRows, _ := h.db.Query(r.Context(), `
		SELECT c.id, c.name, c.brand, c.image_url
		FROM cosmetics c
		JOIN cosmetic_stores cs ON cs.cosmetic_id = c.id
		WHERE cs.store_id = $1 AND c.is_active = TRUE LIMIT 20`, id)
	defer cosRows.Close()
	var cosmetics []map[string]any
	for cosRows.Next() {
		var cid, cname, cbrand string
		var imgURL *string
		cosRows.Scan(&cid, &cname, &cbrand, &imgURL)
		cosmetics = append(cosmetics, map[string]any{"id": cid, "name": cname, "brand": cbrand, "image_url": imgURL})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id": sid, "name": name, "address": addr, "city": city,
		"latitude": lat, "longitude": lng, "website": website, "opening_hours": hours,
		"cosmetics": cosmetics,
	})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
