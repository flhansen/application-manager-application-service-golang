package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"flhansen/application-manager/application-service/src/auth"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareNoToken(t *testing.T) {
	r := httprouter.New()
	mw := AuthMiddleware{SignKey: []byte("supersecretsignkey")}
	srv := http.Server{Addr: "localhost:8000", Handler: r}

	defer srv.Shutdown(context.Background())

	r.GET("/test", mw.Authenticated(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, it's a test.")
	}))

	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		_, err := http.Get("http://localhost:8000/test")

		if err != nil {
			t.Fatal(err)
		}

		return
	case err := <-done:
		t.Fatal(err)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	r := httprouter.New()
	mw := AuthMiddleware{SignKey: []byte("supersecretsignkey")}
	srv := http.Server{Addr: "localhost:8000", Handler: r}

	defer srv.Shutdown(context.Background())

	r.GET("/test", mw.Authenticated(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, it's a test.")
	}))

	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatal(err)
		}

		token, err := auth.GenerateToken(0, "test", jwt.SigningMethodRS256, privateKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		res, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	case err := <-done:
		t.Fatal(err)
	}
}

func TestAuthMiddlewareAuthorized(t *testing.T) {
	r := httprouter.New()
	mw := AuthMiddleware{SignKey: []byte("supersecretsignkey")}
	srv := http.Server{Addr: "localhost:8000", Handler: r}

	defer srv.Shutdown(context.Background())

	r.GET("/test", mw.Authenticated(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, it's a test.")
	}))

	done := make(chan error, 1)
	go func() {
		done <- srv.ListenAndServe()
	}()

	select {
	case <-time.After(500 * time.Millisecond):
		token, err := auth.GenerateToken(0, "test", jwt.SigningMethodHS256, mw.SignKey)
		if err != nil {
			t.Fatal(err)
		}

		client := &http.Client{}

		req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Add("Authorization", token)
		res, err := client.Do(req)

		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusOK, res.StatusCode)
	case err := <-done:
		t.Fatal(err)
	}
}
