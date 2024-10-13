package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lakshay88/real-time-stock/auth"
	"github.com/lakshay88/real-time-stock/config"
	"github.com/lakshay88/real-time-stock/internal/stock/handlers"
	"github.com/lakshay88/real-time-stock/redis"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfiguration("../../config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	// Connect to the database
	// var db database.Database
	// switch cfg.Database.Driver {
	// case "postgres":
	// 	db, err := database.ConnectionToPostgres(cfg.Database)
	// 	if err != nil {
	// 		log.Fatalf("Failed to connect with DB: %v", err)
	// 		return
	// 	}
	// }

	// Set up Redis
	rdb := redis.SetUpRedis()

	// Set up router
	r := mux.NewRouter()

	// Set up stock route, using JWT middleware and passing Redis client and config to handler
	r.Handle("/api/v1/stock/{symbol}", auth.JWTAuthMiddleware(handlers.StockHandler(rdb, *cfg))).Methods("GET")

	// Start the server
	port := ":8081"
	log.Printf("Stock-Service running on port %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// func startGRPC() {
// 	lis, err := net.Listen("tcp", "50051")
// 	if err != nil {
// 		log.Fatalf("Failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()

// 	pb.RegisterStockServiceServer(s, &server{})

// 	log.Println("StockService is running on port :50051")
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("Failed to serve: %v", err)
// 	}
// }

// func (s *server) StreamStockPriceUpdates(req *pb.StockRequest, stream pb.StockService_StreamStockPriceUpdateServer) error {
// 	// Simulate sending stock price updates
// 	for {
// 		// Fetch stock price from your data source or API
// 		price := 150.00 // Replace with actual fetching logic
// 		update := &pb.StockUpdate{
// 			Symbol:    req.Symbol,
// 			Price:     float32(price),
// 			Timestamp: time.Now().Format(time.RFC3339), // Example timestamp
// 		}

// 		if err := stream.Send(update); err != nil {
// 			return err
// 		}

// 		time.Sleep(1 * time.Second) // Adjust as needed
// 	}
// }
