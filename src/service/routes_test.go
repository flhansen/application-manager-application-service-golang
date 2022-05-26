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

func TestRouteDeleteApplication(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	id, err := s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s:%d/api/applications/%d", s.Config.Host, s.Config.Port, id), nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		_, err = s.ApplicationController.GetApplication(id)

		assert.NotNil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteDeleteApplicationRequestError(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	_, err = s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("http://%s:%d/api/applications/foo", s.Config.Host, s.Config.Port), nil)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteUpdateApplication(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	id, err := s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	application, _ := s.ApplicationController.GetApplication(id)

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		application.Commentary = "Test commentary"

		requestBuffer := new(bytes.Buffer)
		if err := json.NewEncoder(requestBuffer).Encode(application); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/applications", s.Config.Host, s.Config.Port), requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteUpdateApplicationInvalidRequestBody(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	_, err = s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		requestBuffer := bytes.NewBufferString(`{ user`)

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/applications", s.Config.Host, s.Config.Port), requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteUpdateApplicationDoesNotExist(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	id, err := s.ApplicationController.InsertApplication(controller.Application{
		UserId:     1,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	application, _ := s.ApplicationController.GetApplication(id)

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		application.Id = -1
		application.Commentary = "Test commentary"

		requestBuffer := new(bytes.Buffer)
		if err := json.NewEncoder(requestBuffer).Encode(application); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/applications", s.Config.Host, s.Config.Port), requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteUpdateApplicationUnauthorized(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.ApplicationController.CreateScheme()
	id, err := s.ApplicationController.InsertApplication(controller.Application{
		UserId:     2,
		WorkTypeId: 1,
		StatusId:   1,
	})

	if err != nil {
		t.Fatal(err)
	}

	application, _ := s.ApplicationController.GetApplication(id)

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		application.Commentary = "Test commentary"

		requestBuffer := new(bytes.Buffer)
		if err := json.NewEncoder(requestBuffer).Encode(application); err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s:%d/api/applications", s.Config.Host, s.Config.Port), requestBuffer)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(1, "testuser", jwt.SigningMethodHS256, s.Config.Jwt.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}
		req.Header.Add("Authorization", token)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetWorkTypes(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.TypesController.CreateScheme()

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/types/worktypes", s.Config.Host, s.Config.Port))
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
		assert.NotNil(t, res["workTypes"])
		assert.Equal(t, 3, len(res["workTypes"].([]interface{})))
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetWorkTypesError(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.TypesController.CreateScheme()

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		if _, err := s.TypesController.Database.Exec(s.TypesController.Context, "DROP TABLE work_type CASCADE"); err != nil {
			t.Fatal(err)
		}

		resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/types/worktypes", s.Config.Host, s.Config.Port))
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetStatuses(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.TypesController.CreateScheme()

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/types/statuses", s.Config.Host, s.Config.Port))
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
		assert.NotNil(t, res["statuses"])
		assert.Equal(t, 3, len(res["statuses"].([]interface{})))
	case err := <-done:
		t.Fatal(err)
	}
}

func TestRouteGetStatusesError(t *testing.T) {
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	s.TypesController.CreateScheme()

	if err != nil {
		t.Fatal(err)
	}

	done := make(chan error)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
		if _, err := s.TypesController.Database.Exec(s.TypesController.Context, "DROP TABLE application_status CASCADE"); err != nil {
			t.Fatal(err)
		}

		resp, err := http.Get(fmt.Sprintf("http://%s:%d/api/types/statuses", s.Config.Host, s.Config.Port))
		if err != nil {
			t.Fatal(err)
		}

		var res map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatal(err)
		}

		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.NotNil(t, res["status"])
		assert.NotNil(t, res["message"])
	case err := <-done:
		t.Fatal(err)
	}
}
