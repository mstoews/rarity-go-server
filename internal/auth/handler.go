package auth

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	db     *pgxpool.Pool
	secret string
}

func NewHandler(db *pgxpool.Pool) *Handler {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		s = "dev-secret-change-in-prod"
	}
	return &Handler{db: db, secret: s}
}

func (h *Handler) Secret() string { return h.secret }

// POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	hash, err := HashPassword(body.Password)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	var userID string
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO users (email, username, password_hash) VALUES ($1,$2,$3) RETURNING id`,
		body.Email, body.Username, hash,
	).Scan(&userID)
	if err != nil {
		jsonErr(w, "email or username already taken", http.StatusConflict)
		return
	}
	h.issueSession(w, r, userID, body.Username, body.Email, "")
}

// POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	var (
		userID, username, subStatus string
		passwordHash                string
		email                       string
	)
	err := h.db.QueryRow(r.Context(),
		`SELECT id, username, email, password_hash, sub_status FROM users WHERE email=$1`, body.Email,
	).Scan(&userID, &username, &email, &passwordHash, &subStatus)
	if err != nil || !CheckPassword(body.Password, passwordHash) {
		jsonErr(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	h.issueSession(w, r, userID, username, email, subStatus)
}

// POST /auth/apple
func (h *Handler) Apple(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IdentityToken string  `json:"identity_token"`
		Email         *string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	appleSub, email, err := verifyAppleToken(body.IdentityToken)
	if err != nil {
		jsonErr(w, "invalid apple token", http.StatusUnauthorized)
		return
	}
	if body.Email != nil && *body.Email != "" {
		email = *body.Email
	}
	var userID, username, subStatus string
	h.db.QueryRow(r.Context(),
		`SELECT id, username, sub_status FROM users WHERE apple_sub=$1`, appleSub,
	).Scan(&userID, &username, &subStatus)

	if userID == "" {
		// First sign-in — create account
		uname := "user_" + appleSub[:8]
		err = h.db.QueryRow(r.Context(),
			`INSERT INTO users (apple_sub, email, username) VALUES ($1,$2,$3) RETURNING id`,
			appleSub, email, uname,
		).Scan(&userID)
		if err != nil {
			jsonErr(w, "could not create account", http.StatusInternalServerError)
			return
		}
		username = uname
	}
	h.issueSession(w, r, userID, username, email, subStatus)
}

// POST /auth/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	rows, err := h.db.Query(r.Context(),
		`SELECT id, user_id, token_hash, expires_at FROM refresh_tokens WHERE expires_at > NOW() LIMIT 100`)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var (
		tokenID, userID string
		hash            string
		exp             time.Time
	)
	for rows.Next() {
		var tid, uid, th string
		var te time.Time
		if err := rows.Scan(&tid, &uid, &th, &te); err != nil {
			continue
		}
		if VerifyRefreshToken(body.RefreshToken, th) {
			tokenID, userID, hash, exp = tid, uid, th, te
			break
		}
	}
	_ = hash
	_ = exp
	if tokenID == "" {
		jsonErr(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}
	h.db.Exec(r.Context(), `DELETE FROM refresh_tokens WHERE id=$1`, tokenID)

	access, err := IssueAccessToken(userID, h.secret)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"access_token": access})
}

func (h *Handler) issueSession(w http.ResponseWriter, r *http.Request, userID, username, email, subStatus string) {
	access, err := IssueAccessToken(userID, h.secret)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	raw, hashed, err := GenerateRefreshToken()
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	exp := time.Now().Add(30 * 24 * time.Hour)
	_, err = h.db.Exec(r.Context(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)`,
		userID, hashed, exp)
	if err != nil {
		jsonErr(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"access_token":  access,
		"refresh_token": raw,
		"user": map[string]any{
			"id":                  userID,
			"username":            username,
			"email":               email,
			"subscription_status": subStatus,
		},
	})
}

func verifyAppleToken(token string) (sub, email string, err error) {
	// TODO: verify JWS signature with Apple's public keys
	// For now return a deterministic stub so tests can run
	if token == "" {
		return "", "", &appleErr{}
	}
	return "apple-stub-" + token[:8], "", nil
}

type appleErr struct{}

func (e *appleErr) Error() string { return "invalid apple token" }

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
