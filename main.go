package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
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

// ==================== TASK 1: GetCountry/:name ====================

func getCountryHandler(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Name parameter is required")
	}

	var country string
	err := db.QueryRow("SELECT Country FROM Person WHERE Name = ?", name).Scan(&country)
	if err == sql.ErrNoRows {
		return c.String(http.StatusNotFound, fmt.Sprintf("Person '%s' not found", name))
	} else if err != nil {
		return c.String(http.StatusInternalServerError, "Database error")
	}

	return c.String(http.StatusOK, country)
}

// ==================== TASK 2: GetCurrentTime/:timezone ====================

func getCurrentTimeHandler(c echo.Context) error {
	timezone := c.Param("timezone")
	if timezone == "" {
		return c.String(http.StatusBadRequest, "Timezone parameter is required")
	}

	// Consume timeapi.io
	url := fmt.Sprintf("https://timeapi.io/api/time/current/zone?timeZone=%s", timezone)
	resp, err := http.Get(url)
	if err != nil {
		return c.String(http.StatusBadGateway, "Failed to call time API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.String(resp.StatusCode, fmt.Sprintf("Time API returned status %d", resp.StatusCode))
	}

	var timeResp TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResp); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse time API response")
	}

	return c.JSON(http.StatusOK, timeResp)
}

// ==================== MAIN ====================

func main() {
	initDB()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	// Task 1 routes
	e.GET("/GetCountry/:name", getCountryHandler)

	// Task 2 routes
	e.GET("/GetCurrentTime/:timezone", getCurrentTimeHandler)

	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"endpoints": "GetCountry/:name, GetCurrentTime/:timezone",
		})
	})

	log.Printf("Server starting on :%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
