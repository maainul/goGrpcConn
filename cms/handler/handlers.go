package handler

import (
	"github.com/gorilla/mux"
)

const (
	index      = "/"
	blog       = "/blogs"
	blogCreate = "/blog/create"
	ErrorPath  = "/error"
)

func handlers(r *mux.Router, s Server) {
	r.HandleFunc(index, s.indexHander).Methods("GET")
	r.HandleFunc(blogCreate, s.blogFormHandler).Methods("GET")
	r.HandleFunc(blogCreate, s.blogCreateHandler).Methods("Post")
}
