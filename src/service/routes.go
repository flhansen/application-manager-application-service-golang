package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (s ApplicationService) handleGetApplications(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// We don't need to check, if userId is not a number, because the
	// authorization middleware does this check for us
	userId, _ := strconv.Atoi(r.Header.Get("userId"))
	applications, _ := s.ApplicationController.GetApplications(userId)

	fmt.Fprint(w, NewApiResponseObject(200, "Fetched all applications", map[string]interface{}{
		"applications": applications,
	}))
}

func handleGetApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleCreateApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleDeleteApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleUpdateApplication(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleGetWorkTypes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}

func handleGetStatuses(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

}
