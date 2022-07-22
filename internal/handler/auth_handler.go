package handler

import (
	"backend/internal/validations"
	"backend/pkg/auth"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type AuthHandler struct {
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewAuthHandler(rd auth.AuthInterface, tk auth.TokenInterface) *AuthHandler {
	return &AuthHandler{
		rd: rd,
		tk: tk,
	}
}

type signUp struct {
	Id       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Admin    bool   `json:"admin" validate:"required"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := validations.SignUpValidation{}
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err := validations.UniversalValidation(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	// TODO: Ganti Ini
	// user := validations.SignUpValidation{}
	user := signUp{}
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if user.Email != "tes" || user.Password != "tes" {
	// 	http.Error(w, "Invalid email or password", http.StatusInternalServerError)
	// 	return
	// }

	ts, err := h.tk.CreateToken(user.Id, user.Admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	saveErr := h.rd.CreateAuth(user.Id, ts)

	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusInternalServerError)
		return
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokens["access_token"])
	w.Write([]byte(tokens["refresh_token"]))
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	metadata, _ := h.tk.ExtractTokenMetadata(r)
	if metadata != nil {
		deleteErr := h.rd.DeleteTokens(metadata)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Write([]byte("Successfully logged out"))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mapToken); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	refreshToken := mapToken["refresh_token"]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		userId, roleOk := claims["user_id"].(string)
		if !roleOk {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check refresh revoked or not
		isRevoked, err := h.rd.CheckRevoked(refreshUuid)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if isRevoked {
			h.rd.RevokeAll(userId)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		} else {
			h.rd.RevokeRefresh(refreshUuid, userId)
		}

		// Using userId to get the user role, temp: admin true
		admin := true

		ts, createErr := h.tk.CreateToken(userId, admin)
		if createErr != nil {
			http.Error(w, createErr.Error(), http.StatusInternalServerError)
			return
		}

		saveErr := h.rd.CreateAuth(userId, ts)

		if saveErr != nil {
			http.Error(w, saveErr.Error(), http.StatusInternalServerError)
			return
		}

		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", "Bearer "+tokens["access_token"])
		w.Write([]byte(tokens["refresh_token"]))
		return
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
}
