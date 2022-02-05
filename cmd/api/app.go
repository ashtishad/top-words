package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const Port = ":8080"

func Start() {

	l := log.New(os.Stdout, "top-ten-api ", log.LstdFlags)

	th := InitHandler(l)

	// created http request multiplexer
	mux := http.NewServeMux()

	// register handlers with routes

	mux.Handle("/topten", th)

	// Server Config
	srv := &http.Server{
		Addr:           Port,
		Handler:        th,
		IdleTimeout:    100 * time.Second,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   2 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			return
		}
	}()

	// graceful shutdown of timeout 1 minute
	gracefulShutdown(srv)
}

// gracefulShutdown wait for interrupt signal to gracefully shut down the server with a timeout of 1 Minute.
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// graceful shutdown
	log.Println("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	log.Println("Server gracefully stopped")
}
