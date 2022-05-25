package service

import (
	"encoding/json"
	"flhansen/application-manager/application-service/src/controller"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type JwtConfig struct {
	SignKey interface{}
}

type ApplicationServiceConfig struct {
	Host     string
	Port     int
	Jwt      JwtConfig
	Database controller.DbConfig
}

type ApplicationService struct {
	Config                ApplicationServiceConfig
	Router                *httprouter.Router
	ApplicationController *controller.ApplicationController
}

func NewApiResponse(status int, message string) string {
	response := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	jsonObj, _ := json.Marshal(response)
	return string(jsonObj)
}

func NewApiResponseObject(status int, message string, moreProps map[string]interface{}) string {
	// Create the default response message
	response := map[string]interface{}{
		"status":  status,
		"message": message,
	}

	// Copy all other properties to the response
	for k, v := range moreProps {
		if _, ok := moreProps[k]; ok {
			response[k] = v
		}
	}

	// Encode JSON object to string
	jsonObj, _ := json.Marshal(response)
	return string(jsonObj)
}

func ApiResponse(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	fmt.Fprint(w, NewApiResponse(code, message))
}

func NewService(config ApplicationServiceConfig) (ApplicationService, error) {
	ac, err := controller.NewApplicationController(config.Database)
	if err != nil {
		return ApplicationService{}, err
	}

	s := ApplicationService{
		Config:                config,
		Router:                httprouter.New(),
		ApplicationController: &ac,
	}

	mw := AuthMiddleware{SignKey: s.Config.Jwt.SignKey}

	// Endpoint: Applications
	s.Router.GET("/api/applications", mw.Authenticated(s.handleGetApplications))
	s.Router.GET("/api/applications/:id", mw.Authenticated(handleGetApplication))
	s.Router.POST("/api/applications", mw.Authenticated(handleCreateApplication))
	s.Router.DELETE("/api/applications/:id", mw.Authenticated(handleDeleteApplication))
	s.Router.PUT("/api/applications/:id", mw.Authenticated(handleUpdateApplication))

	// Endpoint: Types
	s.Router.GET("/api/types/worktypes", handleGetWorkTypes)
	s.Router.GET("/api/types/statuses", handleGetStatuses)

	return s, nil
}

func (s *ApplicationService) Start() error {
	return http.ListenAndServe(fmt.Sprintf("%s:%d", s.Config.Host, s.Config.Port), s.Router)
}
