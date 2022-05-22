package controller

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	Id             int       `db:"id"`
	UserId         int       `db:"user_id"`
	JobTitle       string    `db:"job_title"`
	WorkTypeId     int       `db:"work_type_id"`
	CompanyName    string    `db:"company_name"`
	SubmissionDate time.Time `db:"submission_date"`
	StatusId       int       `db:"status_id"`
	WantedSalary   float32   `db:"wanted_salary"`
	AcceptedSalary float32   `db:"accepted_salary"`
	StartDate      time.Time `db:"start_date"`
	Commentary     string    `db:"commentary"`
}

type ApplicationController struct {
	Database *pgxpool.Pool
	Context  context.Context
}

func NewApplicationController(dbConfig DbConfig) (ApplicationController, error) {
	db, err := pgxpool.Connect(context.Background(), dbConfig.ConnectionString())
	if err != nil {
		return ApplicationController{}, err
	}

	return ApplicationController{
		Database: db,
		Context:  context.Background(),
	}, nil
}

func (c ApplicationController) CreateScheme() {
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
		rows, _ := c.Database.Query(c.Context, query)
		defer rows.Close()
	}
}

func (c ApplicationController) InsertApplication(application Application) (int, error) {
	row := c.Database.QueryRow(c.Context,
		"INSERT INTO application (user_id, job_title, work_type_id, company_name, submission_date, status_id, wanted_salary, accepted_salary, start_date, commentary) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id",
		application.UserId, application.JobTitle, application.WorkTypeId, application.CompanyName, application.SubmissionDate, application.StatusId, application.WantedSalary, application.AcceptedSalary, application.StartDate, application.Commentary)

	id := -1
	err := row.Scan(&id)
	return id, err
}

func (c ApplicationController) GetApplications(userId int) ([]Application, error) {
	rows, err := c.Database.Query(c.Context, "SELECT * FROM application WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application

	for rows.Next() {
		var application Application
		rows.Scan(
			&application.Id, &application.UserId, &application.JobTitle, &application.WorkTypeId, &application.CompanyName,
			&application.SubmissionDate, &application.StatusId, &application.WantedSalary, &application.AcceptedSalary,
			&application.StartDate, &application.Commentary)

		applications = append(applications, application)
	}

	return applications, nil
}
