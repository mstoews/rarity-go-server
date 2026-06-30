package cosmetics

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mstoews/rarity-go-server/internal/middleware"
)

type Handler struct{ db *pgxpool.Pool }

func NewHandler(db *pgxpool.Pool) *Handler { return &Handler{db: db} }

// GET /cosmetics  — free tier gets name/brand/image; paid gets full card with rating
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	categoryID := q.Get("category_id")
	cursor := q.Get("cursor")
	search := q.Get("q")

	userID := middleware.UserIDFrom(r.Context())
	isPaid := h.isSubscribed(r, userID)

	args := []any{}
	sql := `SELECT c.id, c.name, c.brand, c.tagline, c.image_url,
	               cat.id, cat.name,
	               c.avg_rating, c.review_count
	        FROM cosmetics c
	        LEFT JOIN categories cat ON cat.id = c.category_id
	        WHERE c.is_active = TRUE`

	i := 1
	if categoryID != "" {
		sql += ` AND c.category_id = $` + itoa(i)
		args = append(args, categoryID)
		i++
	}
	if search != "" {
		sql += ` AND (c.name ILIKE $` + itoa(i) + ` OR c.brand ILIKE $` + itoa(i) + `)`
		args = append(args, "%"+search+"%")
		i++
	}
	if cursor != "" {
		sql += ` AND c.id > $` + itoa(i)
		args = append(args, cursor)
		i++
	}
	sql += ` ORDER BY c.id LIMIT 21`

	rows, err := h.db.Query(r.Context(), sql, args...)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []map[string]any
	for rows.Next() {
		var id, name, brand string
		var tagline, imageURL, catID, catName *string
		var avgRating float64
		var reviewCount int
		rows.Scan(&id, &name, &brand, &tagline, &imageURL, &catID, &catName, &avgRating, &reviewCount)
		item := map[string]any{"id": id, "name": name, "brand": brand, "tagline": tagline, "image_url": imageURL}
		if catID != nil {
			item["category"] = map[string]any{"id": catID, "name": catName}
		}
		if isPaid {
			item["avg_rating"] = avgRating
			item["review_count"] = reviewCount
		}
		items = append(items, item)
	}

	var nextCursor *string
	if len(items) == 21 {
		last := items[20]["id"].(string)
		nextCursor = &last
		items = items[:20]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"cosmetics": items, "next_cursor": nextCursor})
}

// GET /cosmetics/{id}  — 402 for free users
func (h *Handler) Detail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := middleware.UserIDFrom(r.Context())
	if !h.isSubscribed(r, userID) {
		jsonErr(w, "subscription required", http.StatusPaymentRequired)
		return
	}
	var (
		cid, name, brand               string
		tagline, desc, ingredients     *string
		imageURL                       *string
		images                         []string
		catID, catName                 *string
		avgRating                      float64
		reviewCount                    int
	)
	err := h.db.QueryRow(r.Context(), `
		SELECT c.id, c.name, c.brand, c.tagline, c.description, c.ingredients,
		       c.image_url, c.images,
		       cat.id, cat.name,
		       c.avg_rating, c.review_count
		FROM cosmetics c
		LEFT JOIN categories cat ON cat.id = c.category_id
		WHERE c.id = $1 AND c.is_active = TRUE`, id,
	).Scan(&cid, &name, &brand, &tagline, &desc, &ingredients,
		&imageURL, &images, &catID, &catName, &avgRating, &reviewCount)
	if err != nil {
		jsonErr(w, "not found", http.StatusNotFound)
		return
	}

	// Load stores for this cosmetic
	storeRows, _ := h.db.Query(r.Context(), `
		SELECT s.id, s.name, s.address, s.city, s.latitude, s.longitude, cs.in_stock, cs.notes
		FROM stores s
		JOIN cosmetic_stores cs ON cs.store_id = s.id
		WHERE cs.cosmetic_id = $1`, id)
	defer storeRows.Close()
	var stores []map[string]any
	for storeRows.Next() {
		var sid, sname string
		var addr, city *string
		var lat, lng *float64
		var inStock *bool
		var notes *string
		storeRows.Scan(&sid, &sname, &addr, &city, &lat, &lng, &inStock, &notes)
		stores = append(stores, map[string]any{
			"id": sid, "name": sname, "address": addr, "city": city,
			"latitude": lat, "longitude": lng, "in_stock": inStock, "notes": notes,
		})
	}

	item := map[string]any{
		"id": cid, "name": name, "brand": brand, "tagline": tagline,
		"description": desc, "ingredients": ingredients, "image_url": imageURL, "images": images,
		"avg_rating": avgRating, "review_count": reviewCount, "stores": stores,
	}
	if catID != nil {
		item["category"] = map[string]any{"id": catID, "name": catName}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) isSubscribed(r *http.Request, userID string) bool {
	if userID == "" {
		return false
	}
	var status string
	h.db.QueryRow(r.Context(),
		`SELECT sub_status FROM users WHERE id=$1 AND (sub_status='active') AND (sub_expires_at IS NULL OR sub_expires_at > NOW())`,
		userID).Scan(&status)
	return status == "active"
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func itoa(i int) string {
	return string(rune('0' + i))
}
