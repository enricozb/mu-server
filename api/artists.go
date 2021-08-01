package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) artists(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.lib.Artists()); err != nil {
		http.Error(w, fmt.Errorf("encode: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
