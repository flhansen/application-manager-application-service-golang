package controller

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var controller ApplicationController
var testApplication Application

func TestMain(m *testing.M) {
	os.Exit(runAllTests(m))
}

func runAllTests(m *testing.M) int {
	controller, _ = NewApplicationController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	testApplication = Application{
		UserId:         1,
		JobTitle:       "test job",
		WorkTypeId:     1,
		CompanyName:    "test company",
		SubmissionDate: time.Now(),
		StatusId:       1,
		WantedSalary:   0,
		AcceptedSalary: 0,
		StartDate:      time.Now(),
		Commentary:     "empty",
	}

	return m.Run()
}

func TestNewApplicationControllerConnectionError(t *testing.T) {
	controller, err := NewApplicationController(DbConfig{
		Host:     "localhost",
		Port:     1234,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	assert.Nil(t, controller.Context)
	assert.Nil(t, controller.Database)
	assert.NotNil(t, err)
}

func TestNewApplicationController(t *testing.T) {
	controller, err := NewApplicationController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	assert.NotNil(t, controller)
	assert.Nil(t, err)
}

func TestCreateScheme(t *testing.T) {
	controller.CreateScheme()
}

func TestInsertApplication(t *testing.T) {
	controller.CreateScheme()

	_, err := controller.InsertApplication(Application{
		UserId:      0,
		JobTitle:    "test",
		CompanyName: "test",
		WorkTypeId:  1,
		StatusId:    1,
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetApplications(t *testing.T) {
	controller.CreateScheme()
	userId := 1

	_, err := controller.InsertApplication(Application{UserId: userId, WorkTypeId: 1, StatusId: 1})
	if err != nil {
		t.Fatal(err)
	}

	var numberApplications int
	row := controller.Database.QueryRow(controller.Context, "SELECT count(*) FROM application")
	if err = row.Scan(&numberApplications); err != nil {
		t.Fatal(err)
	}

	applications, err := controller.GetApplications(userId)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, numberApplications, len(applications))
}

func TestGetApplicationsError(t *testing.T) {
	controller.CreateScheme()

	rows, err := controller.Database.Query(controller.Context, "DROP TABLE application")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	applications, err := controller.GetApplications(-1)

	assert.Nil(t, applications)
	assert.NotNil(t, err)
}

func TestGetApplication(t *testing.T) {
	controller.CreateScheme()

	id, err := controller.InsertApplication(testApplication)
	if err != nil {
		t.Fatal(err)
	}

	application, err := controller.GetApplication(id)

	assert.Nil(t, err)
	assert.Equal(t, testApplication.UserId, application.UserId)
	assert.Equal(t, testApplication.JobTitle, application.JobTitle)
	assert.Equal(t, testApplication.CompanyName, application.CompanyName)
}

func TestGetApplicationDoesNotExist(t *testing.T) {
	controller.CreateScheme()
	_, err := controller.GetApplication(1)

	assert.NotNil(t, err)
}

func TestDeleteApplication(t *testing.T) {
	controller.CreateScheme()

	id, err := controller.InsertApplication(testApplication)
	if err != nil {
		t.Fatal(err)
	}

	controller.DeleteApplication(id)
	_, err = controller.GetApplication(id)

	assert.NotNil(t, err)
}

func TestUpdateApplication(t *testing.T) {
	controller.CreateScheme()

	id, err := controller.InsertApplication(testApplication)
	if err != nil {
		t.Fatal(err)
	}

	updatedApplication := testApplication
	updatedApplication.Id = id
	updatedApplication.JobTitle = "another job title"
	controller.UpdateApplication(updatedApplication)

	application, err := controller.GetApplication(id)

	assert.Nil(t, err)
	assert.Equal(t, "another job title", application.JobTitle)
}
