package api

import (
	"encoding/json"
	"fmt"
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
func (a *API) songCover(w http.ResponseWriter, r *http.Request) {}

// TODO: add ->AAC conversion
//        stream the output of the ffmpeg conversion
// TODO: add cover extraction
