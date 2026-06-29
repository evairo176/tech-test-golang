package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// ==================== MODELS ====================

type Person struct {
	Name    string `json:"name"`
	Country string `json:"country"`
}

type TimeResponse struct {
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Day          int    `json:"day"`
	Hour         int    `json:"hour"`
	Minute       int    `json:"minute"`
	Seconds      int    `json:"seconds"`
	MilliSeconds int    `json:"milliSeconds"`
	DateTime     string `json:"dateTime"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	TimeZone     string `json:"timeZone"`
	DayOfWeek    string `json:"dayOfWeek"`
	DstActive    bool   `json:"dstActive"`
}

// ==================== DATABASE ====================

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./persons.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Person (
		Name    TEXT PRIMARY KEY,
		Country TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Insert data (only if empty)
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM Person").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		insertSQL := `
		INSERT INTO Person (Name, Country) VALUES
			('Adam', 'Kuala Lumpur'),
			('John', 'Singapore'),
			('Henry', 'Singapore'),
			('Dominic', 'Thailand');`
		_, err = db.Exec(insertSQL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Seeded Person table with initial data")
	}
}

// ==================== TASK 1: GetCountry/{Name} ====================

func getCountryHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		// Fallback for Go < 1.22 style or manual extraction
		// Try extracting from URL path
		path := r.URL.Path
		prefix := "/GetCountry/"
		if len(path) > len(prefix) {
			name = path[len(prefix):]
		}
	}

	if name == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}

	var country string
	err := db.QueryRow("SELECT Country FROM Person WHERE Name = ?", name).Scan(&country)
	if err == sql.ErrNoRows {
		http.Error(w, fmt.Sprintf("Person '%s' not found", name), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(country))
}

// ==================== TASK 2: GetCurrentTime/{Timezone} ====================

func getCurrentTimeHandler(w http.ResponseWriter, r *http.Request) {
	timezone := r.PathValue("timezone")
	if timezone == "" {
		path := r.URL.Path
		prefix := "/GetCurrentTime/"
		if len(path) > len(prefix) {
			timezone = path[len(prefix):]
		}
	}

	if timezone == "" {
		http.Error(w, "Timezone parameter is required", http.StatusBadRequest)
		return
	}

	// Consume timeapi.io
	url := fmt.Sprintf("https://timeapi.io/api/time/current/zone?timeZone=%s", timezone)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to call time API", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Time API returned status %d", resp.StatusCode), resp.StatusCode)
		return
	}

	var timeResp TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResp); err != nil {
		http.Error(w, "Failed to parse time API response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(timeResp)
}

// ==================== MAIN ====================

func main() {
	initDB()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Task 1 routes
	mux.HandleFunc("/GetCountry/{name}", getCountryHandler)
	mux.HandleFunc("/GetCountry/", getCountryHandler)

	// Task 2 routes
	mux.HandleFunc("/GetCurrentTime/{timezone}", getCurrentTimeHandler)
	mux.HandleFunc("/GetCurrentTime/", getCurrentTimeHandler)

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"endpoints": "GetCountry/{name}, GetCurrentTime/{timezone}",
		})
	})

	log.Printf("Server starting on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
