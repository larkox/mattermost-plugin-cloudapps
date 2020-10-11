package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/server/apps"
	"github.com/mattermost/mattermost-plugin-apps/server/constants"
	"github.com/mattermost/mattermost-plugin-apps/server/utils/httputils"
)

const (
	SubscribePath = "/subscribe"
)

type restapi struct {
	apps *apps.Service
}

func Init(router *mux.Router, apps *apps.Service) {
	a := restapi{
		apps: apps,
	}

	subrouter := router.PathPrefix(constants.APIPath).Subrouter()
	subrouter.HandleFunc(SubscribePath, a.handleSubscribe).Methods("POST", "DELETE")

	subrouter.HandleFunc("/locations", checkAuthorized(a.handleLocations)).Methods("GET")
}

func checkAuthorized(f func(http.ResponseWriter, *http.Request, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		actingUserID := req.Header.Get("Mattermost-User-Id")
		if actingUserID == "" {
			httputils.WriteUnauthorizedError(w, errors.New("not authorized"))
			return
		}

		f(w, req, actingUserID)
	}
}
