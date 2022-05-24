package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTypesController(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	assert.Nil(t, err)
	assert.NotNil(t, controller.Context)
	assert.NotNil(t, controller.Database)
}

func TestNewTypesControllerDatabaseConnectionError(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     1234,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	assert.NotNil(t, err)
	assert.Nil(t, controller.Context)
	assert.Nil(t, controller.Database)
}

func TestGetWorkTypes(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	controller.CreateScheme()
	workTypes, err := controller.GetWorkTypes()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(workTypes))
}

func TestGetWorkTypesNoRelationError(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	controller.CreateScheme()
	if _, err = controller.Database.Exec(controller.Context, "DROP TABLE work_type CASCADE"); err != nil {
		t.Fatal(err)
	}

	_, err = controller.GetWorkTypes()

	assert.NotNil(t, err)
}

func TestGetStatuses(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	controller.CreateScheme()
	statuses, err := controller.GetStatuses()

	assert.Nil(t, err)
	assert.Equal(t, 3, len(statuses))
}

func TestGetStatusesNoRelationError(t *testing.T) {
	controller, err := NewTypesController(DbConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "test",
		Password: "test",
		Database: "test",
	})

	if err != nil {
		t.Fatal(err)
	}

	controller.CreateScheme()
	if _, err = controller.Database.Exec(controller.Context, "DROP TABLE application_status CASCADE"); err != nil {
		t.Fatal(err)
	}

	_, err = controller.GetStatuses()

	assert.NotNil(t, err)
}
