package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	// Load environment variables
	err := godotenv.Load("../.env") // Sesuaikan dengan lokasi file .env
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Ambil konfigurasi database dari environment
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Buat connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	fmt.Println("Connecting to database with:", connStr) // Log koneksi database

	// Hubungkan ke database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Database not responding:", err)
	}
	fmt.Println("✅ Connected to database successfully!")
}

type Message struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func getMessages(w http.ResponseWriter, r *http.Request) {
	fmt.Println("📩 Received request: GET /messages") // Log request masuk

	rows, err := db.Query("SELECT id, text FROM messages")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("❌ Error fetching messages:", err)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Text); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println("❌ Error scanning message:", err)
			return
		}
		messages = append(messages, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
	fmt.Println("✅ Successfully sent response:", messages) // Log response yang dikirim
}

// Middleware untuk mengizinkan CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	initDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/messages", getMessages).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

	// Tambahkan middleware CORS
	handler := enableCORS(r)

	fmt.Println("🚀 Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
