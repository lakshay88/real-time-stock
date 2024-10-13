package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/lakshay88/real-time-stock/auth"
	"github.com/lakshay88/real-time-stock/config"
	"github.com/lakshay88/real-time-stock/database"
	"github.com/lakshay88/real-time-stock/internal/user/handlers"
)

const (
	restPort = ":8080"
	grpcPort = ":50051"
)

func main() {
	// Creating to goroutines to finish
	cfg, err := config.LoadConfiguration("../../config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	var db database.Database

	switch cfg.Database.Driver {
	case "postgres":
		db, err = database.ConnectionToPostgres(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		startGRPCServer()
	}()

	go func() {
		defer wg.Done()
		startRESTServer(db)
	}()

	wg.Wait()
}

func startGRPCServer() {
}

func startRESTServer(db database.Database) {

	// Registring Router
	r := mux.NewRouter()

	// Routs
	r.HandleFunc("/api/v1/createUser", handlers.CreateUserHandler(db)).Methods("POST")
	r.HandleFunc("/api/v1/login", handlers.LoginUser(db)).Methods("POST")
	r.Handle("/api/v1/getAllUser", auth.JWTAuthMiddleware(http.HandlerFunc(handlers.GetAllUser(db)))).Methods("GET")

	log.Printf("REST Serves Started running on port: %s", restPort)

	// Starting Serves
	if err := http.ListenAndServe(restPort, r); err != nil {
		log.Fatal("Failed to start REST Serves")
	}
}
