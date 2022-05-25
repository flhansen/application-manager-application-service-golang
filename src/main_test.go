package main

import (
	"flhansen/application-manager/application-service/src/service"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestRunApplication(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	config := service.ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: service.JwtConfig{
			SignKey: "supersecretsignkey",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	configFileName := filepath.Join(os.TempDir(), "test_config.yml")
	if err := ioutil.WriteFile(configFileName, data, 0777); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(configFileName)

	os.Args[1] = configFileName

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		return
	case exitCode := <-done:
		t.Fatalf("The application terminated with %d\n", exitCode)
	}
}

func TestRunApplicationNoSuchFile(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	configFileName := filepath.Join(os.TempDir(), "test_config.yml")
	os.Args[1] = configFileName

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatal("The application should not run without a configuration file\n")
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	}
}

func TestRunApplicationInvalidConfig(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	config := "invalid file content"

	configFileName := filepath.Join(os.TempDir(), "test_config.yml")
	if err := ioutil.WriteFile(configFileName, []byte(config), 0777); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(configFileName)

	os.Args[1] = configFileName

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatal("The application should not run without a valid configuration file\n")
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	}
}

func TestRunApplicationStartError(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	config := service.ApplicationServiceConfig{
		Host: "localhost",
		Port: -1,
		Jwt: service.JwtConfig{
			SignKey: "supersecretsignkey",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}

	configFileName := filepath.Join(os.TempDir(), "test_config.yml")
	if err := ioutil.WriteFile(configFileName, data, 0777); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(configFileName)

	os.Args[1] = configFileName

	done := make(chan int)
	go func() {
		done <- runApplication()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		t.Fatal("The application should not run with invalid configuration\n")
	case exitCode := <-done:
		assert.Equal(t, 1, exitCode)
	}
}
