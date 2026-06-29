package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"tech-test-golang/internal/model"
)

type PersonRepository struct {
	db *sql.DB
}

func NewPersonRepository(db *sql.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

// GetAll returns all persons from the database
func (r *PersonRepository) GetAll() ([]model.Person, error) {
	rows, err := r.db.Query("SELECT Name, Country FROM Person ORDER BY Name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := []model.Person{}
	for rows.Next() {
		var p model.Person
		if err := rows.Scan(&p.Name, &p.Country); err != nil {
			return nil, err
		}
		persons = append(persons, p)
	}
	return persons, nil
}

// GetCountryByName returns the country for a given person name (case-insensitive)
func (r *PersonRepository) GetCountryByName(name string) (string, error) {
	var country string
	err := r.db.QueryRow("SELECT Country FROM Person WHERE LOWER(Name) = LOWER(?)", name).Scan(&country)
	return country, err
}

// Create inserts a new person into the database
func (r *PersonRepository) Create(p model.Person) (int64, error) {
	result, err := r.db.Exec("INSERT INTO Person (Name, Country) VALUES (?, ?)", p.Name, p.Country)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return 0, fmt.Errorf("person '%s' already exists", p.Name)
		}
		return 0, err
	}
	return result.RowsAffected()
}

// Delete removes a person by name (case-insensitive)
func (r *PersonRepository) Delete(name string) (int64, error) {
	result, err := r.db.Exec("DELETE FROM Person WHERE LOWER(Name) = LOWER(?)", name)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count returns the total number of persons
func (r *PersonRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM Person").Scan(&count)
	return count, err
}
