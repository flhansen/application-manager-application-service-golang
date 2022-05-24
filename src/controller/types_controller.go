package controller

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type WorkType struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type ApplicationStatus struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type TypesController struct {
	Database *pgxpool.Pool
	Context  context.Context
}

func NewTypesController(dbConfig DbConfig) (TypesController, error) {
	controller := TypesController{
		Context: context.Background(),
	}

	db, err := pgxpool.Connect(controller.Context, dbConfig.ConnectionString())
	if err != nil {
		return TypesController{}, err
	}

	controller.Database = db
	return controller, nil
}

// TODO: Code duplication (see ApplicationController.CreateScheme). Refactor it,
// so both structs share the same logic.
func (c TypesController) CreateScheme() {
	queries := []string{
		`DROP TABLE IF EXISTS work_type CASCADE`,
		`CREATE TABLE work_type (
				id SERIAL PRIMARY KEY NOT NULL,
				name VARCHAR(255) NOT NULL)`,
		`DROP TABLE IF EXISTS application_status CASCADE`,
		`CREATE TABLE application_status (
				id SERIAL PRIMARY KEY NOT NULL,
				name VARCHAR(255) NOT NULL)`,
		`DROP TABLE IF EXISTS application CASCADE`,
		`CREATE TABLE application (
				id SERIAL PRIMARY KEY NOT NULL,
				user_id INTEGER NOT NULL,
				job_title VARCHAR(255) NOT NULL,
				work_type_id INTEGER NOT NULL,
				company_name VARCHAR(255) NOT NULL,
				submission_date DATE,
				status_id INTEGER NOT NULL,
				wanted_salary REAL,
				accepted_salary REAL,
				start_date DATE,
				commentary VARCHAR(500),

				FOREIGN KEY (work_type_id) REFERENCES work_type (id)
					ON DELETE SET DEFAULT,
				FOREIGN KEY (status_id) REFERENCES application_status (id)
					ON DELETE SET DEFAULT)`,
		`INSERT INTO work_type (name) VALUES ('Remote')`,
		`INSERT INTO work_type (name) VALUES ('OnSite')`,
		`INSERT INTO work_type (name) VALUES ('Hybrid')`,
		`INSERT INTO application_status (name) VALUES ('Accepted')`,
		`INSERT INTO application_status (name) VALUES ('Pending')`,
		`INSERT INTO application_status (name) VALUES ('Declined')`,
	}

	for _, query := range queries {
		c.Database.Exec(c.Context, query)
	}
}

func (c TypesController) GetWorkTypes() ([]WorkType, error) {
	rows, err := c.Database.Query(c.Context, "SELECT * FROM work_type")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var workTypes []WorkType
	for rows.Next() {
		var workType WorkType
		rows.Scan(&workType.Id, &workType.Name)
		workTypes = append(workTypes, workType)
	}

	return workTypes, nil
}

func (c TypesController) GetStatuses() ([]ApplicationStatus, error) {
	rows, err := c.Database.Query(c.Context, "SELECT * FROM application_status")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var statuses []ApplicationStatus
	for rows.Next() {
		var status ApplicationStatus
		rows.Scan(&status.Id, &status.Name)
		statuses = append(statuses, status)
	}

	return statuses, nil
}
