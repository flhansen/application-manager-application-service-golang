package service

import (
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

	srv := &http.Server{Addr: fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), Handler: s.Router}
	defer srv.Shutdown(context.Background())

	if err != nil {
		t.Fatal(err)
	}

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
	case err := <-done:
		t.Fatal(err)
	}
}
