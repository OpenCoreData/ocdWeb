package dx

import (
	"log"
	"net/http"
	// "net/url"
)

// This is a place holder for a generic openSearch endpoint..
// it needs to get an opensearch formated call and route to the correct search..
// since this is routing, do I want to use the gorilla routing capacity then too?
func Redirection(w http.ResponseWriter, r *http.Request) {

	// call mongo and lookup the redirection to use...

	log.Printf("dx:  %s", r.URL.Path)
	// toto  update with regex /ref/  for /id/
	http.Redirect(w, r, "/ref/"+r.URL.Path, 303) //  need a hnader for the redirects and put this in that handler)
}
