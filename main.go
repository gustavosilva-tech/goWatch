package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func allowGet(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	return true
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	writeJSON(w, map[string]string{
		"status":  "ok",
		"service": "gowatch",
	})
}

func checkHandler(client *http.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowGet(w, r) {
			return
		}

		writeJSON(w, checkURL(client, "http://127.0.0.1:8080/health"))
	}
}

func checksHandler(client *http.Client, targets []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowGet(w, r) {
			return
		}

		writeJSON(w, checkURLs(client, targets, 2))
	}
}

func statusHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !allowGet(w, r) {
			return
		}

		writeJSON(w, store.Latest())
	}
}

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	targets := []string{
		"http://127.0.0.1:8080/health",
		"http://127.0.0.1:8080/not-found",
	}

	store := NewStore()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/check", checkHandler(client))
	mux.HandleFunc("/checks", checksHandler(client, targets))
	mux.HandleFunc("/status", statusHandler(store))

	server := &http.Server{
		Addr:              "127.0.0.1:8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,

	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	log.Println("GoWatch running at http://localhost:8080")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		runScheduler(ctx, client, targets, 2, 10*time.Second, store)
	}()

	<-ctx.Done()
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	wg.Wait()
	log.Println("GoWatch stopped")
}
