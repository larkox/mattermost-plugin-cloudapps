package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-plugin-apps/server/store"
)

type SubscribeResponse struct {
	Error  string            `json:"error,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

func (a *restapi) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	var err error
	actingUserID := ""
	// logMessage := ""
	status := http.StatusInternalServerError

	defer func() {
		resp := SubscribeResponse{}
		if err != nil {
			resp.Error = errors.Wrap(err, "failed to subscribe").Error()
			// logMessage = "Error: " + resp.Error
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(resp)
	}()

	actingUserID = r.Header.Get("Mattermost-User-ID")
	if actingUserID == "" {
		err = errors.New("user not logged in")
		status = http.StatusUnauthorized
		return
	}
	// TODO check for sysadmin

	var sub store.Subscription
	if err = json.NewDecoder(r.Body).Decode(&sub); err != nil {
		status = http.StatusUnauthorized
		return
	}

	// TODO replace with an appropriate API-level call that would validate,
	// deduplicate, etc.
	switch r.Method {
	case http.MethodPost:
		err = a.apps.Store.StoreSub(&sub)
	case http.MethodDelete:
		err = a.apps.Store.DeleteSub(&sub)
	default:
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
