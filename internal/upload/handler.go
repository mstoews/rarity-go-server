package upload

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Handler struct {
	presign *s3.PresignClient
	bucket  string
}

func NewHandler() (*Handler, error) {
	bucket := os.Getenv("S3_BUCKET")
	if bucket == "" {
		bucket = "rarity-uploads"
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	presign := s3.NewPresignClient(client)
	return &Handler{presign: presign, bucket: bucket}, nil
}

// POST /upload/presigned-url
func (h *Handler) PresignedURL(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Filename string `json:"filename"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		jsonErr(w, "bad request", http.StatusBadRequest)
		return
	}
	key := fmt.Sprintf("reviews/%d-%s", time.Now().UnixNano(), body.Filename)
	req, err := h.presign.PresignPutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.bucket),
		Key:         aws.String(key),
		ContentType: aws.String("image/jpeg"),
	}, s3.WithPresignExpires(15*time.Minute))
	if err != nil {
		jsonErr(w, "could not generate upload URL", http.StatusInternalServerError)
		return
	}
	imageURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", h.bucket, key)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"upload_url": req.URL, "image_url": imageURL})
}

func jsonErr(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
