package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
)

type SignInRequest struct {
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		WriteJson(w, SignInResponse{Error: "Авторизация не настроена"})
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJson(w, SignInResponse{Error: "Ошибка чтения JSON"})
		return
	}

	if req.Password != pass {
		WriteJson(w, SignInResponse{Error: "Неверный пароль"})
		return
	}

	hash := sha256.Sum256([]byte(pass))
	token := hex.EncodeToString(hash[:])

	WriteJson(w, SignInResponse{
		Token: token,
	})
}
