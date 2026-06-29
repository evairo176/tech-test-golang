package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"tech-test-golang/internal/model"
	"tech-test-golang/internal/repository"
	"tech-test-golang/internal/response"

	"github.com/labstack/echo/v4"
)

type PersonHandler struct {
	repo *repository.PersonRepository
}

func NewPersonHandler(repo *repository.PersonRepository) *PersonHandler {
	return &PersonHandler{repo: repo}
}

// GetAllPersons handles GET /persons
func (h *PersonHandler) GetAllPersons(c echo.Context) error {
	persons, err := h.repo.GetAll()
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to query persons")
	}
	return response.Success(c, http.StatusOK, fmt.Sprintf("Found %d persons", len(persons)), persons)
}

// GetCountryByName handles GET /getcountry/:name
func (h *PersonHandler) GetCountryByName(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return response.Error(c, http.StatusBadRequest, "Name parameter is required")
	}

	country, err := h.repo.GetCountryByName(name)
	if err == sql.ErrNoRows {
		return response.Error(c, http.StatusNotFound, fmt.Sprintf("Person '%s' not found", name))
	} else if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Database error")
	}

	return response.Success(c, http.StatusOK, fmt.Sprintf("Country for '%s' found", name), map[string]string{
		"name":    name,
		"country": country,
	})
}

// CreatePerson handles POST /person
func (h *PersonHandler) CreatePerson(c echo.Context) error {
	p := new(model.Person)
	if err := c.Bind(p); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}
	if p.Name == "" || p.Country == "" {
		return response.Error(c, http.StatusBadRequest, "Name and Country are required")
	}

	rowsAffected, err := h.repo.Create(*p)
	if err != nil {
		// Check for duplicate
		if fmt.Sprintf("%v", err) == fmt.Sprintf("person '%s' already exists", p.Name) {
			return response.Error(c, http.StatusConflict, err.Error())
		}
		return response.Error(c, http.StatusInternalServerError, "Failed to insert person")
	}

	return response.Success(c, http.StatusCreated, "Person created successfully", map[string]interface{}{
		"name":         p.Name,
		"country":      p.Country,
		"rowsAffected": rowsAffected,
	})
}

// DeletePerson handles DELETE /person/:name
func (h *PersonHandler) DeletePerson(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return response.Error(c, http.StatusBadRequest, "Name parameter is required")
	}

	rowsAffected, err := h.repo.Delete(name)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to delete person")
	}
	if rowsAffected == 0 {
		return response.Error(c, http.StatusNotFound, fmt.Sprintf("Person '%s' not found", name))
	}

	return response.Success(c, http.StatusOK, fmt.Sprintf("Person '%s' deleted", name), nil)
}

// RootHandler handles GET /
func RootHandler(c echo.Context) error {
	return response.Success(c, http.StatusOK, "Tech Test Golang API", map[string]interface{}{
		"endpoints": []map[string]string{
			{"method": "POST", "path": "/person", "description": "Create a new person (body: {name, country})"},
			{"method": "GET", "path": "/persons", "description": "Get all persons"},
			{"method": "GET", "path": "/GetCountry/{name}", "description": "Get country by person name"},
			{"method": "DELETE", "path": "/person/{name}", "description": "Delete person by name"},
			{"method": "GET", "path": "/GetCurrentTime/{timezone}", "description": "Get current time by timezone"},
		},
	})
}
func GetCurrentTime(c echo.Context) error {
	timezone := c.Param("timezone")
	if timezone == "" {
		return response.Error(c, http.StatusBadRequest, "Timezone parameter is required")
	}

	decoded, err := url.QueryUnescape(timezone)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid timezone format")
	}
	timezone = decoded

	requestURL := fmt.Sprintf("https://timeapi.io/api/time/current/zone?timeZone=%s", url.QueryEscape(timezone))
	resp, err := http.Get(requestURL)
	if err != nil {
		return response.Error(c, http.StatusBadGateway, "Failed to call time API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.Error(c, resp.StatusCode, fmt.Sprintf("Time API returned status %d", resp.StatusCode))
	}

	var timeResp model.TimeResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeResp); err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to parse time API response")
	}

	return response.Success(c, http.StatusOK, fmt.Sprintf("Current time for '%s'", timezone), timeResp)
}
