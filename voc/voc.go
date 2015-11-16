package voc

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func VocJanus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Show Janus voc entry version %s of %s", vars["version"], vars["term"])
	fmt.Fprintf(w, "%s", "template for this term coming")
}

func VocCore(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Show Core voc entry version %s of %s", vars["version"], vars["term"])
	fmt.Fprintf(w, "%s", "template for this term coming")
}
