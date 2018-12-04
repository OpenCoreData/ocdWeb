package dx

import (
	"log"
	"net/http"
	"strings"
)

// Redirection handler
// I had functions for RDFRedirection and Expedition
//  Why did I have that?   Was there a plan for the
// existence of two other functions?
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

// ProvRedirection is a function to provide prov on a resource
// It is method GET.   Rewrite and send back for routing...
func ProvRedirection(w http.ResponseWriter, r *http.Request) {
	log.Printf("dx:  %s", r.URL.Path)
	newPath := strings.Replace(r.URL.Path, "/id/", "/rdf/", 1)
	http.Redirect(w, r, newPath, 303)
}

// PingbackRedirection is a function to handle POST calls to this
// endpoint and load the results to a store.  I do not need
// rewrite..  this could be removed..  the call in main.go would
// simply route to a handle function
func PingbackRedirection(w http.ResponseWriter, r *http.Request) {
	log.Printf("dx:  %s", r.URL.Path)
	newPath := strings.Replace(r.URL.Path, "/id/", "/rdf/", 1)
	http.Redirect(w, r, newPath, 303)
}
