package handlers

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/lakshay88/real-time-stock/config"
	"github.com/lakshay88/real-time-stock/internal/stock/thirdparty"
)

func StockHandler(rdb *redis.Client, cfg config.AppConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		symbol := mux.Vars(r)["symbol"]

		// Check cache
		cacheKey := "stock_" + symbol
		cachedResponse, err := rdb.Get(r.Context(), cacheKey).Result()
		if err == nil {
			log.Printf("Cache hit for symbol: %s", symbol)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedResponse)) // Serve from cache
			return
		}

		// Cache miss, call third-party API
		log.Printf("Cache miss for symbol: %s, fetching from third-party API", symbol)
		stockRequest := thirdparty.SetUpStockAPI(cfg)

		responseDataBytes, err := stockRequest.GetStockData(symbol)
		if err != nil {
			log.Fatalf("Error calling third-party API: %v", err)
			http.Error(w, "Failed to fetch stock data", http.StatusInternalServerError)
			return
		}

		sendBroadcastRequest(responseDataBytes)

		// Cache the response
		rdb.Set(r.Context(), cacheKey, string(responseDataBytes), 10*time.Second)
		log.Println("Response cached with 5 seconds TTL")

		// Send the API response
		w.Header().Set("Content-Type", "application/json")
		w.Write(responseDataBytes)
	}
}

func sendBroadcastRequest(data []byte) {
	url := "http://localhost:8082/api/v1/broadcast"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send broadcast request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to broadcast message: %v", resp.Status)
	}
}
