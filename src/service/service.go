package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type JwtConfig struct {
	SignKey interface{}
}

type ApplicationServiceConfig struct {
	Host string
	Port int
	Jwt  JwtConfig
}

type ApplicationService struct {
	Config ApplicationServiceConfig
	Router *httprouter.Router
}

func NewApiResponse(status int, message string) string {
	response := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	jsonObj, _ := json.Marshal(response)
	return string(jsonObj)
}

func ApiResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	fmt.Fprint(w, NewApiResponse(code, message))
}

func NewService(config ApplicationServiceConfig) ApplicationService {
	s := ApplicationService{
		Config: config,
		Router: httprouter.New(),
	}

	mw := AuthMiddleware{SignKey: s.Config.Jwt.SignKey}

	// Endpoint: Applications
	s.Router.GET("/api/applications", mw.Authenticated(handleGetApplications))
	s.Router.GET("/api/applications/:id", mw.Authenticated(handleGetApplication))
	s.Router.POST("/api/applications", mw.Authenticated(handleCreateApplication))
	s.Router.DELETE("/api/applications/:id", mw.Authenticated(handleDeleteApplication))
	s.Router.PUT("/api/applications/:id", mw.Authenticated(handleUpdateApplication))

	// Endpoint: Types
	s.Router.GET("/api/types/worktypes", handleGetWorkTypes)
	s.Router.GET("/api/types/statuses", handleGetStatuses)

	return s
}

func (s ApplicationService) Start() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), s.Router)
}
