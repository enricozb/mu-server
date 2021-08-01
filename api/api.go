package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/enricozb/mu-server/library"
)

type API struct {
	lib *library.Library
}

func New(lib *library.Library) *API {
	return &API{lib: lib}
}

func (a *API) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/songs", a.songs).Methods("GET")
	r.HandleFunc("/songs/{id}", a.song).Methods("GET")
	r.HandleFunc("/songs/{id}/cover", a.songCover).Methods("GET")

	r.HandleFunc("/albums", a.albums).Methods("GET")
	r.HandleFunc("/albums/{id}/cover", a.albumCover).Methods("GET")

	r.HandleFunc("/artists", a.artists).Methods("GET")

	http.Handle("/", r)

	fmt.Printf("serving %d items...\n", a.lib.Size())
	return http.ListenAndServe(":4000", cors.Default().Handler(r))
}
