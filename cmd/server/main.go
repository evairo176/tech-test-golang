package main

import (
	"log"
	"strings"

	"tech-test-golang/config"
	"tech-test-golang/internal/handler"
	"tech-test-golang/internal/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	// Load config
	cfg := config.Load()

	// Init database with migrations & seeders
	db := config.InitDB(cfg)
	defer db.Close()

	// Init repository
	personRepo := repository.NewPersonRepository(db)

	// Init handlers
	personHandler := handler.NewPersonHandler(personRepo)

	// Setup Echo
	e := echo.New()

	// Case-insensitive routing middleware
	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().URL.Path = strings.ToLower(c.Request().URL.Path)
			return next(c)
		}
	})

	// API group
	api := e.Group("/api")

	// Task 1 — Person CRUD
	api.POST("/person", personHandler.CreatePerson)
	api.GET("/persons", personHandler.GetAllPersons)
	api.GET("/getcountry/:name", personHandler.GetCountryByName)
	api.DELETE("/person/:name", personHandler.DeletePerson)

	// Task 2 — Time API
	api.GET("/getcurrenttime/:timezone", handler.GetCurrentTime)

	// Root endpoint
	e.GET("/", func(c echo.Context) error {
		return handler.RootHandler(c)
	})

	log.Printf("Server starting on :%s\n", cfg.Port)
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
