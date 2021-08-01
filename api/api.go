package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type API struct {
}

func (a *API) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/songs", a.songs).Methods("GET")
	r.HandleFunc("/songs/{id}", a.song).Methods("GET")
	r.HandleFunc("/songs/{id}/cover", a.songCover).Methods("GET")
	r.HandleFunc("/songs/{id}/metadata", a.songMetadata).Methods("GET")

	r.HandleFunc("/albums", a.albums).Methods("GET")
	r.HandleFunc("/albums/{id}/cover", a.albumCover).Methods("GET")
	r.HandleFunc("/albums/{id}/songs", a.albumSongs).Methods("GET")

	r.HandleFunc("/artists", a.artists).Methods("GET")
	r.HandleFunc("/artists/{id}/songs", a.artistSongs).Methods("GET")

	http.Handle("/", r)

	return http.ListenAndServe(":4000", cors.Default().Handler(r))
}
