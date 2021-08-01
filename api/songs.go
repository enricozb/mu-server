package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *API) songs(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(a.lib.Songs()); err != nil {
		http.Error(w, fmt.Errorf("encode: %v", err).Error(), http.StatusInternalServerError)
		return
	}
}

func (a *API) song(w http.ResponseWriter, r *http.Request)      {}
func (a *API) songCover(w http.ResponseWriter, r *http.Request) {}

// TODO: add ->AAC conversion
//        stream the output of the ffmpeg conversion
// TODO: add cover extraction
