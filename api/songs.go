package api

import "net/http"

func (a *API) songs(w http.ResponseWriter, r *http.Request)        {}
func (a *API) song(w http.ResponseWriter, r *http.Request)         {}
func (a *API) songCover(w http.ResponseWriter, r *http.Request)    {}
func (a *API) songMetadata(w http.ResponseWriter, r *http.Request) {}
