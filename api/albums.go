package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) albums(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.lib.Albums()); err != nil {
		http.Error(w, fmt.Errorf("encode: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
func (a *API) albumCover(w http.ResponseWriter, r *http.Request) {}
