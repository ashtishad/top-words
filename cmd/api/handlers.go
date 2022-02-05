package api

import (
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
		h.serveTopTen(w, r)
	} else {
		h.l.Println("Invalid request method: ", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// core handler to serve top ten words
func (h *TopTen) serveTopTen(w http.ResponseWriter, r *http.Request) {
	defer service.Init() // reset after use
	h.l.Println("Endpoint hit: /top-ten")
	w.Header().Set("Content-Type", "application/json")

	// receive text from json request body
	data := r.Body
	defer data.Close()

	// create a new dto.TopTen object
	text, err := service.FromJSON(data)

	// send text to the top ten handler
	c := service.Init()
	resp, err := c.GetTopTenWords(text)
	if err != nil {
		h.l.Println("Error while getting top ten words : ", err.AsMessage())
		http.Error(w, err.AsMessage(), http.StatusInternalServerError)
		return
	}

	h.l.Println(resp) // test

	// send response to the client

	err = service.ToJSON(w, resp)
	if err != nil {
		h.l.Println("Error while sending response to json : ", err.AsMessage())
		http.Error(w, err.AsMessage(), http.StatusInternalServerError)
		return
	}
}
