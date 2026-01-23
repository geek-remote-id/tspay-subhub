package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Category represents a category in the cashier system
type Category struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	DeletedAt   *timestamppb.Timestamp `json:"deleted_at"`
}

// Response represents the standardized API response format
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// WriteJSON is a helper to write JSON responses
func WriteJSON(w http.ResponseWriter, status int, res Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using system environment variables")
	}

	// connect to DB
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	defer db.Close()

	// check connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	fmt.Println("Successfully connected to database!")

	fmt.Println("Successfully connected to database!")

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, Response{
			Status:  "success",
			Message: "API Running",
		})
	})

	// POST localhost:8080/api/category
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT id, name, description, deleted_at FROM category")
			if err != nil {
				WriteJSON(w, http.StatusInternalServerError, Response{
					Status:  "failed",
					Message: "Failed to fetch categories: " + err.Error(),
				})
				return
			}
			defer rows.Close()

			var categories []Category
			for rows.Next() {
				var c Category
				var deletedAt sql.NullTime
				if err := rows.Scan(&c.ID, &c.Name, &c.Description, &deletedAt); err != nil {
					WriteJSON(w, http.StatusInternalServerError, Response{
						Status:  "failed",
						Message: "Failed to scan category: " + err.Error(),
					})
					return
				}
				if deletedAt.Valid {
					c.DeletedAt = timestamppb.New(deletedAt.Time)
				}
				categories = append(categories, c)
			}

			WriteJSON(w, http.StatusOK, Response{
				Status:  "success",
				Message: "Categories retrieved successfully",
				Data:    categories,
			})

		case "POST":
			// baca data dari request
			var categoryBaru Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				WriteJSON(w, http.StatusBadRequest, Response{
					Status:  "failed",
					Message: "Invalid request body",
				})
				return
			}

			// simpan ke database
			var deletedAt sql.NullTime
			err = db.QueryRow(
				"INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id, deleted_at",
				categoryBaru.Name, categoryBaru.Description,
			).Scan(&categoryBaru.ID, &deletedAt)

			if err != nil {
				WriteJSON(w, http.StatusInternalServerError, Response{
					Status:  "failed",
					Message: "Failed to save category: " + err.Error(),
				})
				return
			}

			if deletedAt.Valid {
				categoryBaru.DeletedAt = timestamppb.New(deletedAt.Time)
			}

			WriteJSON(w, http.StatusCreated, Response{
				Status:  "success",
				Message: "Category created successfully",
				Data:    categoryBaru,
			})
		default:
			WriteJSON(w, http.StatusMethodNotAllowed, Response{
				Status:  "failed",
				Message: "Method not allowed",
			})
		}
	})

	fmt.Println("Server running on http://localhost:8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error running server:", err)
	}
}
