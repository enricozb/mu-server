package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *API) albums(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.lib.Albums()); err != nil {
		http.Error(w, fmt.Errorf("encode: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
func (a *API) albumCover(w http.ResponseWriter, r *http.Request) {
	album := mux.Vars(r)["id"]
	data, err := a.lib.AlbumCover(album)

	if err != nil {
		http.Error(w, fmt.Errorf("album cover: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, data); err != nil {
		http.Error(w, fmt.Errorf("copy: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
