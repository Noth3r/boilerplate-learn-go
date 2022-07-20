package handlers

import (
	"backend/services/auth"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type authHandler struct {
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewAuthHandler(rd auth.AuthInterface, tk auth.TokenInterface) *authHandler {
	return &authHandler{
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

func (h *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user := signUp{}
	if err := decoder.Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Email != "tes" || user.Password != "tes" {
		http.Error(w, "Invalid email or password", http.StatusInternalServerError)
		return
	}

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

func (h *authHandler) SignOut(w http.ResponseWriter, r *http.Request) {
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

func (h *authHandler) Refresh(w http.ResponseWriter, r *http.Request) {
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

		// Using userId to get the user role, temp: admin true
		admin := true

		delErr := h.rd.DeleteRefresh(refreshUuid)
		if (delErr) != nil {
			http.Error(w, delErr.Error(), http.StatusInternalServerError)
			return
		}

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

// type AuthHandler struct {
// 	server *s.Server
// }

// func NewAuthHandler(server *s.Server) *AuthHandler {
// 	return &AuthHandler{server: server}
// }

// type signUp struct {
// 	Email    string `json:"email" validate:"required,email"`
// 	Password string `json:"password" validate:"required"`
// }

// func (a *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	payload := signUp{}
// 	if err := decoder.Decode(&payload); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if payload.Email != "tes" || payload.Password != "tes" {
// 		http.Error(w, "Invalid email or password", http.StatusInternalServerError)
// 		return
// 	}

// 	data, err := generateToken()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Header().Set("Authorization", "Bearer "+data["refresh_token"])
// 	w.Write([]byte(data["access_token"]))
// }

// func generateToken() (map[string]string, error) {
// 	token := jwt.New(jwt.SigningMethodHS256)
// 	claims := token.Claims.(jwt.MapClaims)
// 	claims["email"] = "tes"
// 	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
// 	claims["admin"] = false
// 	t, err := token.SignedString([]byte("secret"))

// 	refreshToken := jwt.New(jwt.SigningMethodHS256)
// 	rtClaims := refreshToken.Claims.(jwt.MapClaims)
// 	rtClaims["sub"] = 1
// 	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
// 	rt, err := refreshToken.SignedString([]byte("secret"))

// 	if err != nil {
// 		return nil, err
// 	}

// 	// data, err := json.Marshal()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return map[string]string{"access_token": t, "refresh_token": rt}, nil
// }
