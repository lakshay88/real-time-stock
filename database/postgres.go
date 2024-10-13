package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lakshay88/real-time-stock/config"
	"github.com/lakshay88/real-time-stock/internal/user/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	connection *sql.DB
}

func (p *PostgresDB) CreateUserIfNotExist() error {
	query := `Create Table IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL
	)`

	_, err := p.connection.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create user table - %s", err)
		return err
	}
	return nil
}

func (p *PostgresDB) CreateUser(user models.User) error {

	err := p.CreateUserIfNotExist()
	if err != nil {
		return err
	}

	query := `INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4)`
	_, err = p.connection.Exec(query, user.Username, user.Email, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute insert query: %w", err)
	}

	return nil
}

func (p *PostgresDB) GetUserList(user []models.User) ([]models.User, error) {
	query := `SELECT id, username, email, created_at FROM users`
	rows, err := p.connection.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil

}

func (p *PostgresDB) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	// Use QueryRow when expecting a single result
	query := `SELECT id, username, email, password, created_at FROM users WHERE email=$1`
	row := p.connection.QueryRow(query, email)

	// Scan the result into the user object
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		// Handle "no rows found" error separately
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user with email %s not found", email)
		}
		return user, err
	}
	return user, nil
}

func ConnectionToPostgres(cfg config.DatabaseConfig) (Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	connection, err := sql.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Set connection limits
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)
	connection.SetConnMaxLifetime(5 * time.Minute)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{connection: connection}, nil
}

// func (p *PostgresDB) Close(cfg config.DatabaseConfig) {
// 	p.connection.Close()
// }
