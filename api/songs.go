package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/enricozb/mu-server/convert"
)

func (a *API) songs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.lib.Songs()); err != nil {
		http.Error(w, fmt.Errorf("encode: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) song(w http.ResponseWriter, r *http.Request) {
	song := mux.Vars(r)["id"]
	path, err := a.lib.Abs(song)
	if err != nil {
		http.Error(w, fmt.Errorf("abs: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	if err := convert.Convert(path, w); err != nil {
		http.Error(w, fmt.Errorf("convert: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
func (a *API) songCover(w http.ResponseWriter, r *http.Request) {
	song := mux.Vars(r)["id"]
	data, err := a.lib.SongCover(song)

	if err != nil {
		http.Error(w, fmt.Errorf("song cover: %v", err).Error(), http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(w, data); err != nil {
		http.Error(w, fmt.Errorf("copy: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}
