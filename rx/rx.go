package rx

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RenderWithProvHeader displays the RDF resource and adds a prov pingback entry
func RenderWithProvHeader(w http.ResponseWriter, r *http.Request) {
	linkProv := fmt.Sprintf("<http://opencoredata.org/id/%s/provenance>; rel=\"http://www.w3.org/ns/prov#has_provenance\"", r.URL.Path[1:])
	linkPB := fmt.Sprintf("<http://opencoredata.org/rdf/%s/pingback>; rel=\"http://www.w3.org/ns/prov#pingbck\"", r.URL.Path[1:])
	w.Header().Add("Link", linkProv)
	w.Header().Add("Link", linkPB)
	w.Header().Set("Content-type", "text/plain")
	fmt.Println(r.URL.Path[1:])
	http.ServeFile(w, r, fmt.Sprintf("./static/%s", r.URL.Path[1:]))
}

// RenderWithProv shows the prov of a resource
func RenderWithProv(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Cool prov here"))
}

// ProvPingback Handles the PROV pingback on a resource
func ProvPingback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	// do something with the POST data
	// likely convert to triples and write to some end point...
	fmt.Printf("Prov for %s\n", r.URL.Path[1:])
	fmt.Println(string(body))

	w.Write([]byte("Thanks for your contribution"))
}
