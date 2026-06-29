package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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

// ==================== STANDARD API RESPONSE ====================

type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(c echo.Context, code int, message string, data interface{}) error {
	return c.JSON(code, APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, APIResponse{
		Status:  "error",
		Code:    code,
		Message: message,
	})
}

// ==================== MIDDLEWARE ====================

func CaseInsensitiveRouting(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		c.Request().URL.Path = strings.ToLower(path)
		return next(c)
	}
}

// ==================== DATABASE ====================

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./persons.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS Person (
		Name    TEXT PRIMARY KEY,
		Country TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

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

// ==================== TASK 1: Person CRUD ====================

// POST /person — Insert new person
func createPersonHandler(c echo.Context) error {
	p := new(Person)
	if err := c.Bind(p); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}
	if p.Name == "" || p.Country == "" {
		return ErrorResponse(c, http.StatusBadRequest, "Name and Country are required")
	}

	result, err := db.Exec("INSERT INTO Person (Name, Country) VALUES (?, ?)", p.Name, p.Country)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrorResponse(c, http.StatusConflict, fmt.Sprintf("Person '%s' already exists", p.Name))
		}
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to insert person")
	}

	rowsAffected, _ := result.RowsAffected()
	return SuccessResponse(c, http.StatusCreated, "Person created successfully", map[string]interface{}{
		"name":         p.Name,
		"country":      p.Country,
		"rowsAffected": rowsAffected,
	})
}

// GET /persons — Select all persons
func getAllPersonsHandler(c echo.Context) error {
	rows, err := db.Query("SELECT Name, Country FROM Person ORDER BY Name")
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to query persons")
	}
	defer rows.Close()

	persons := []Person{}
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.Name, &p.Country); err != nil {
			return ErrorResponse(c, http.StatusInternalServerError, "Failed to scan person")
		}
		persons = append(persons, p)
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("Found %d persons", len(persons)), persons)
}

// GET /getcountry/:name — Get country by person name
func getCountryHandler(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return ErrorResponse(c, http.StatusBadRequest, "Name parameter is required")
	}

	var country string
	err := db.QueryRow("SELECT Country FROM Person WHERE LOWER(Name) = LOWER(?)", name).Scan(&country)
	if err == sql.ErrNoRows {
		return ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Person '%s' not found", name))
	} else if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Database error")
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("Country for '%s' found", name), map[string]string{
		"name":    name,
		"country": country,
	})
}

// DELETE /person/:name — Delete person by name
func deletePersonHandler(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return ErrorResponse(c, http.StatusBadRequest, "Name parameter is required")
	}

	result, err := db.Exec("DELETE FROM Person WHERE LOWER(Name) = LOWER(?)", name)
	if err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to delete person")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Person '%s' not found", name))
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("Person '%s' deleted", name), nil)
}

// ==================== TASK 2: GetCurrentTime/:timezone ====================

func getCurrentTimeHandler(c echo.Context) error {
	timezone := c.Param("timezone")
	if timezone == "" {
		return ErrorResponse(c, http.StatusBadRequest, "Timezone parameter is required")
	}

	decoded, err := url.QueryUnescape(timezone)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, "Invalid timezone format")
	}
	timezone = decoded

	requestURL := fmt.Sprintf("https://timeapi.io/api/time/current/zone?timeZone=%s", url.QueryEscape(timezone))
	resp, err := http.Get(requestURL)
	if err != nil {
		return ErrorResponse(c, http.StatusBadGateway, "Failed to call time API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrorResponse(c, resp.StatusCode, fmt.Sprintf("Time API returned status %d", resp.StatusCode))
	}

	var timeResp TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResp); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, "Failed to parse time API response")
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("Current time for '%s'", timezone), timeResp)
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
	e.Pre(CaseInsensitiveRouting)

	// Task 1 — Person CRUD
	e.POST("/person", createPersonHandler)
	e.GET("/persons", getAllPersonsHandler)
	e.GET("/getcountry/:name", getCountryHandler)
	e.DELETE("/person/:name", deletePersonHandler)

	// Task 2 — Time API
	e.GET("/getcurrenttime/:timezone", getCurrentTimeHandler)

	// Root
	e.GET("/", func(c echo.Context) error {
		return SuccessResponse(c, http.StatusOK, "Tech Test Golang API", map[string]interface{}{
			"endpoints": []map[string]string{
				{"method": "POST", "path": "/person", "description": "Create a new person (body: {name, country})"},
				{"method": "GET", "path": "/persons", "description": "Get all persons"},
				{"method": "GET", "path": "/GetCountry/{name}", "description": "Get country by person name"},
				{"method": "DELETE", "path": "/person/{name}", "description": "Delete person by name"},
				{"method": "GET", "path": "/GetCurrentTime/{timezone}", "description": "Get current time by timezone"},
			},
		})
	})

	log.Printf("Server starting on :%s\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
