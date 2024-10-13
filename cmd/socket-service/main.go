package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lakshay88/real-time-stock/auth"
	"github.com/lakshay88/real-time-stock/internal/socket/handlers"
)

func main() {

	r := mux.NewRouter()

	r.Handle("/api/v1/socketSubscribe", auth.JWTAuthMiddleware(handlers.SocketHandler()))
	r.Handle("/api/v1/broadcast", handlers.BroadcastHandler()).Methods("POST")

	port := ":8082"
	log.Printf("Stock-Service running on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
