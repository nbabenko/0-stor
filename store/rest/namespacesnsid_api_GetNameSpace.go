package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zero-os/0-stor/store/manager"
	"github.com/zero-os/0-stor/store/stats"
)

// GetNameSpace is the handler for GET /namespaces/{nsid}
// Get detail view about namespace
func (api NamespacesAPI) GetNameSpace(w http.ResponseWriter, r *http.Request) {
	label := mux.Vars(r)["nsid"]

	// FIXME: uncomment when the IYO middleware will create the namespace object in
	// db automaticly
	mgr := manager.NewNamespaceManager(api.db)

	// ns, err := mgr.Get(label)
	// if err != nil {
	// 	if err == db.ErrNotFound {
	// 		jsonError(w, err, http.StatusNotFound)
	// 	} else {
	// 		jsonError(w, err, http.StatusInternalServerError)
	// 	}
	// 	return
	// }
	count, err := mgr.Count(label)
	if err != nil {
		jsonError(w, err, http.StatusInternalServerError)
		return
	}
	read, write := stats.Rate(label)
	respBody := Namespace{
		Label: label,
		Stats: NamespaceStat{
			NrObjects:           int64(count),
			ReadRequestPerHour:  read,
			WriteRequestPerHour: write,
			// SpaceAvailable: TODO
			// SpaceUsed: TODO
		},
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&respBody)
}