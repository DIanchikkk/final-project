package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJson(w, map[string]string{
				"error": "Authentication required",
			})
			return
		}

		hash := sha256.Sum256([]byte(pass))
		expected := hex.EncodeToString(hash[:])

		if cookie.Value != expected {
			w.WriteHeader(http.StatusUnauthorized)
			WriteJson(w, map[string]string{
				"error": "Authentication required",
			})
			return
		}

		next(w, r)
	})
}
