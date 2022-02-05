package api

import (
	"encoding/json"
	"github.com/ashtishad/top-words/internal/dto"
	"github.com/ashtishad/top-words/pkg/service"
	"log"
	"net/http"
)

type TopTen struct {
	l *log.Logger
}

func InitHandler(l *log.Logger) *TopTen {
	return &TopTen{l}
}

func (h *TopTen) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		h.l.Println("Handling post request on /top-ten")
		h.serveTopTen(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// core handler to serve top ten words
func (h *TopTen) serveTopTen(w http.ResponseWriter, r *http.Request) {
	defer service.Init() // reset after use
	h.l.Println("Endpoint hit: /top-ten")

	// receive text from json request body
	read := r.Body
	defer read.Close()

	// create a new dto.TopTen object
	var text dto.TextRequestDto
	if err := json.NewDecoder(read).Decode(&text); err != nil {
		h.l.Println("Error while decoding json")
		http.Error(w, "Error while decoding json", http.StatusBadRequest)
		return
	}

	// send text to the top ten handler
	c := service.Init() // reset for next use
	resp, err := c.GetTopTenWords(text)
	if err != nil {
		h.l.Println("Error while getting top ten words")
		http.Error(w, "Error while getting top ten words", http.StatusInternalServerError)
		return
	}

	h.l.Println(resp) // test

	// send response to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.l.Println("Error while encoding json")
		http.Error(w, "Error while encoding json", http.StatusInternalServerError)
		return
	}
}
