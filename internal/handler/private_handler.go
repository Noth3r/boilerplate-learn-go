package handler

import (
	"net/http"
)

type PrivateHandler struct {
}

func (p *PrivateHandler) Tes(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Private routes"))
}
