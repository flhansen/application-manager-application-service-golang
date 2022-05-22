package controller

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var controller ApplicationController

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
