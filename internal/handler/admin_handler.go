package handler

import (
	"net/http"
)

type AdminHandler struct {
}

func (a *AdminHandler) Tes(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello, Admin!"))
}
