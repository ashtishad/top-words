package api

import (
	"log"
	"net/http"
	"os"
)

const Port = ":8080"

func Start() {

	l := log.New(os.Stdout, "top-ten-api ", log.LstdFlags)
	th := InitHandler(l)
	mux := http.NewServeMux()

	mux.Handle("/topten", th)

	http.ListenAndServe(Port, mux)
}
