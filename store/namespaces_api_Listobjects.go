package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger/badger"
	"github.com/gorilla/mux"
)

// Listobjects is the handler for GET /namespaces/{nsid}/objects
// List keys of the namespaces
func (api NamespacesAPI) Listobjects(w http.ResponseWriter, r *http.Request) {
	var respBody []Object

	// Pagination
	pageParam := r.FormValue("page")
	per_pageParam := r.FormValue("per_page")

	if pageParam == "" {
		pageParam = "1"
	}

	if per_pageParam == "" {
		per_pageParam = strconv.Itoa(api.config.Pagination.PageSize)
	}

	page, err := strconv.Atoi(pageParam)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	per_page, err := strconv.Atoi(per_pageParam)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	nsid := mux.Vars(r)["nsid"]

	prefix := fmt.Sprintf("%s:", nsid)

	opt := badger.DefaultIteratorOptions
	opt.PrefetchSize = api.config.Iterator.PreFetchSize

	it := api.db.store.NewIterator(opt)
	defer it.Close()

	startingIndex := (page-1)*per_page + 1
	counter := 0 // Number of objects encountered
	resultsCount := per_page

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		key := string(item.Key()[:])
		/* Skip namespaces and objects not belonging to intended
		   namespace
		*/
		if !strings.Contains(key, prefix) {
			continue
		}

		// Found a namespace
		counter++

		// Skip this object if its index < intended startingIndex
		if counter < startingIndex {
			continue
		}

		value := item.Value()

		var file = &File{}
		object := file.ToObject(value, key)

		respBody = append(respBody, *object)

		if len(respBody) == resultsCount {
			break
		}
	}

	// return empty list if no results
	if len(respBody) == 0 {
		respBody = []Object{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&respBody)
}