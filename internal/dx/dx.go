package dx

import (
	"log"
	"net/http"
	"strings"
)

// Redirection handler
func Redirection(w http.ResponseWriter, r *http.Request) {
	log.Printf("dx:  %s", r.URL.Path)
	newPath := strings.Replace(r.URL.Path, "/id/", "/doc/", 1)
	http.Redirect(w, r, newPath, 303)
}

// RDFRedirection  to get RDF graphs in /rdf/*
func RDFRedirection(w http.ResponseWriter, r *http.Request) {
	log.Printf("dx:  %s", r.URL.Path)
	newPath := strings.Replace(r.URL.Path, "/id/", "/rdf/", 1)
	http.Redirect(w, r, newPath, 303)
}
