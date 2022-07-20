package handlers

import (
	s "backend/server"
	"net/http"
)

type PrivateHandler struct {
	server *s.Server
}

func NewPrivateHandler(server *s.Server) *PrivateHandler {
	return &PrivateHandler{server: server}
}

func (p *PrivateHandler) Tes(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Private routes"))
}

// func (p *PrivateHandler) Token(w http.ResponseWriter, r *http.Request) {
// 	decoder := json.NewDecoder(r.Body)
// 	tokenReq := struct {
// 		RefreshToken string `json:"refresh_token"`
// 	}{}
// 	if err := decoder.Decode(&tokenReq); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return []byte("secret"), nil
// 	})

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		if int(claims["sub"].(float64)) == 1 {
// 			newTokenPair, err := generateToken()
// 			if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Header().Set("Authorization", "Bearer "+newTokenPair["refresh_token"])
// 			w.Write([]byte(newTokenPair["access_token"]))
// 			return
// 		}
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	http.Error(w, err.Error(), http.StatusInternalServerError)
// }
