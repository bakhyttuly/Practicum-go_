package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bookstore/middleware"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id query param required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID < 1 {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	claims := jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		http.Error(w, "token generation error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": signed})
}
