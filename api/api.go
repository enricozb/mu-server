package api

import (
	"fmt"
	"net/http"
	"net/url"

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

func log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if path, err := url.PathUnescape(r.RequestURI); err != nil {
			fmt.Printf("%s: %s\n", r.Method, r.RequestURI)
		} else {
			fmt.Printf("%s: %s\n", r.Method, path)
		}
		next.ServeHTTP(w, r)
	})
}

func (a *API) Run() error {
	r := mux.NewRouter()

	r.Use(log)

	r.HandleFunc("/songs/{id:.+}/cover", a.songCover).Methods("GET")
	r.HandleFunc("/songs/{id:.+}", a.song).Methods("GET")
	r.HandleFunc("/songs", a.songs).Methods("GET")

	r.HandleFunc("/albums/{id:.*}/cover", a.albumCover).Methods("GET")
	r.HandleFunc("/albums", a.albums).Methods("GET")

	r.HandleFunc("/artists", a.artists).Methods("GET")

	r.NewRoute().HandlerFunc(http.NotFound)

	http.Handle("/", r)

	fmt.Printf("serving %d items...\n", a.lib.Size())
	return http.ListenAndServe(":4000", cors.Default().Handler(r))
}
