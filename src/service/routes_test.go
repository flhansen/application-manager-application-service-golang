package service

import (
	"bytes"
	"context"
	"encoding/json"
	"flhansen/application-manager/application-service/src/auth"
	"flhansen/application-manager/application-service/src/controller"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestRouteGetApplications(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/applications", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetApplicationInvalidId(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/applications/foo", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetApplicationNotFound(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/applications/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetApplicationNotAuthorized(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()
	s.ApplicationController.InsertApplication(controller.Application{
		UserId:     2,
		WorkTypeId: 1,
		StatusId:   1,
	})

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/applications/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetApplication(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()
	s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/applications/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
		assert.NotNil(t, res["application"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteCreateApplicationRequestBodyError(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		requestBuffer := new(bytes.Buffer)
		bytes.NewBufferString(`{ userId: `)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/api/applications", requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteCreateApplicationInsertionError(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		application := controller.Application{
			JobTitle: "test",
		}

		requestBuffer := new(bytes.Buffer)
		if err := json.NewEncoder(requestBuffer).Encode(application); err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/api/applications", requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteCreateApplication(t *testing.T) {
	s, err := NewService(ApplicationServiceConfig{
		Host: "localhost",
		Port: 8000,
		Jwt: JwtConfig{
			SignKey: []byte("supersecretsignkey"),
		},
		Database: controller.DbConfig{
			Host:     "localhost",
			Port:     5432,
			Username: "test",
			Password: "test",
			Database: "test",
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	s.ApplicationController.CreateScheme()

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		application := controller.Application{
			JobTitle:   "test",
			WorkTypeId: 1,
			StatusId:   1,
		}

		requestBuffer := new(bytes.Buffer)
		if err := json.NewEncoder(requestBuffer).Encode(application); err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, "http://localhost:8000/api/applications", requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&res)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
		assert.NotNil(t, res["application"])
	case err := <-done:
		t.Fatal(err)
	}
}
