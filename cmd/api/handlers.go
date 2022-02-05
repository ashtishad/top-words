package api

import (
	"log"
	"net/http"
)

type topTenHandler struct {
	l *log.Logger
}

func InitHandler(l *log.Logger) *topTenHandler {
	return &topTenHandler{l}
}

func (h *topTenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.l.Println("Entered into topTenHandler")
	// not implemented

}
