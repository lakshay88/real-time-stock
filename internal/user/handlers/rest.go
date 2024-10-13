package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lakshay88/real-time-stock/auth"
	"github.com/lakshay88/real-time-stock/database"
	"github.com/lakshay88/real-time-stock/internal/user/models"
	"github.com/lakshay88/real-time-stock/utils"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUserHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateUserRequest

		// Decoding Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Reqest Payload", http.StatusBadRequest)
			return
		}

		log.Printf("creating User using Rest")

		// Generate a new UUID for the user ID
		userID := uuid.New().String()
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		user := models.User{
			ID:        userID,
			Username:  req.Username,
			Email:     req.Email,
			Password:  hashedPassword,
			CreatedAt: time.Now(),
		}

		if err := db.CreateUser(user); err != nil {
			return
		}

	}
}

func GetAllUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// write Login code
		users, err := db.GetUserList([]models.User{})
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	}
}

func LoginUser(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginUserRequest

		// Decoding Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		user, err := db.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWTToken(user)
		if err != nil {
			http.Error(w, "Failed to generate JWT token", http.StatusUnauthorized)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})

	}
}
