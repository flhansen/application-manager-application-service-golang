package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewApiResponse(t *testing.T) {
	response := NewApiResponse(200, "Hello, it's a test.")

	var res map[string]interface{}
	err := json.Unmarshal([]byte(response), &res)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, res["status"])
	assert.NotNil(t, res["message"])
	assert.Equal(t, 200.0, res["status"])
	assert.Equal(t, "Hello, it's a test.", res["message"])
}

func TestApiResponse(t *testing.T) {
	var buffer bytes.Buffer
	ApiResponse(&buffer, "Hello, it's a test.", http.StatusOK)

	var obj map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &obj)

	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, obj["status"])
	assert.NotNil(t, obj["message"])
	assert.Equal(t, 200.0, obj["status"])
	assert.Equal(t, "Hello, it's a test.", obj["message"])
}

func TestNewService(t *testing.T) {
	s := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8080,
		Jwt: JwtConfig{
			SignKey: "supersecretsignkey",
		},
	})

	assert.NotNil(t, s)
}

func TestServiceStart(t *testing.T) {
	s := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8080,
		Jwt: JwtConfig{
			SignKey: "supersecretsignkey",
		},
	})

	done := make(chan error, 1)
	go func() {
		done <- s.Start()
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		return
	case err := <-done:
		t.Fatal(err)
	}
}
