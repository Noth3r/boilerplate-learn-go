package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PublicHandler struct {
}

func (p *PublicHandler) Hello(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func (p *PublicHandler) Post(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	payload := struct {
		Tes string `json:"tes"`
	}{}
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(payload.Tes)
	w.Write([]byte("Hello, World!"))
}
