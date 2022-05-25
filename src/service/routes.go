package service

import (
	"encoding/json"
	"flhansen/application-manager/application-service/src/controller"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (s ApplicationService) handleGetApplications(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// We don't need to check, if userId is not a number, because the
	// authorization middleware does this check for us
	userId, _ := strconv.Atoi(p.ByName("userId"))
	applications, _ := s.ApplicationController.GetApplications(userId)

	fmt.Fprint(w, NewApiResponseObject(200, "Fetched all applications", map[string]interface{}{
		"applications": applications,
	}))
}

func (s ApplicationService) handleGetApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	applicationId, err := strconv.Atoi(p.ByName("id"))
	if err != nil {
		ApiResponse(w, "Error while parsing the application id", http.StatusBadRequest)
		return
	}

	application, err := s.ApplicationController.GetApplication(applicationId)
	if err != nil {
		ApiResponse(w, "This application does not exist", http.StatusBadRequest)
		return
	}

	userId, _ := strconv.Atoi(p.ByName("userId"))
	if application.UserId != userId {
		ApiResponse(w, "You are not allowed to get information about this application", http.StatusUnauthorized)
		return
	}

	fmt.Fprint(w, NewApiResponseObject(http.StatusOK, "Fetched application", map[string]interface{}{
		"application": application,
	}))
}

func (s ApplicationService) handleCreateApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var applicationRequest controller.Application
	if err := json.NewDecoder(r.Body).Decode(&applicationRequest); err != nil {
		ApiResponse(w, "Could not create application", http.StatusInternalServerError)
		return
	}

	// Make sure that the application is assigned to the requesting user.
	applicationRequest.UserId, _ = strconv.Atoi(p.ByName("userId"))

	id, err := s.ApplicationController.InsertApplication(applicationRequest)
	if err != nil {
		ApiResponse(w, "Could not create application", http.StatusInternalServerError)
		return
	}

	newApplication, _ := s.ApplicationController.GetApplication(id)
	fmt.Fprint(w, NewApiResponseObject(http.StatusOK, "Application created", map[string]interface{}{
		"application": newApplication,
	}))
}

func handleDeleteApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleUpdateApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleGetWorkTypes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleGetStatuses(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
