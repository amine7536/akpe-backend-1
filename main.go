package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/backend1?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Init schema
	initDB()

	http.HandleFunc("/items", handleItems)
	http.HandleFunc("/health", handleHealth)

	log.Printf("backend-1 listening on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func initDB() {
	schema := `
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL
	);
	INSERT INTO items (name)
	SELECT name FROM (VALUES ('alpha'), ('bravo'), ('charlie')) AS v(name)
	WHERE NOT EXISTS (SELECT 1 FROM items);
	`
	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	log.Println("Database schema initialized")
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name FROM items ORDER BY id")
	if err != nil {
		http.Error(w, fmt.Sprintf("query error: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, fmt.Sprintf("scan error: %v", err), http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"service": "backend-1",
		"items":   items,
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if err := db.Ping(); err != nil {
		http.Error(w, "unhealthy", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
