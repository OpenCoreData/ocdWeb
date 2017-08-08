package rx

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RenderWithProvHeader displays the RDF resource and adds a prov pingback entry
func RenderWithProvHeader(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path[5:])                                                                                                       // Note the HACK in the next line..  this is just an ALPHA..  (even this this sucks)
	linkProv := fmt.Sprintf("<http://%s/id/%s/provenance>; rel=\"http://www.w3.org/ns/prov#has_provenance\"", r.Host, r.URL.Path[5:]) // use r.Host so we don't hardcode in
	linkPB := fmt.Sprintf("<http://%s/rdf/%s/pingback>; rel=\"http://www.w3.org/ns/prov#pingbck\"", r.Host, r.URL.Path[5:])
	w.Header().Add("Link", linkProv)
	w.Header().Add("Link", linkPB)
	w.Header().Set("Content-type", "text/plain")
	fmt.Println(r.URL.Path[1:])
	http.ServeFile(w, r, fmt.Sprintf("./static/%s", r.URL.Path[1:]))
}

// RenderWithProv shows the prov of a resource
// right now it just hist getProvRecord which returns a generic same for all
// record (since I have no prov data stood up now beyond testing stuff)
func RenderWithProv(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getProvRecord()))
}

// ProvPingback Handles the PROV pingback on a resource
func ProvPingback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	fmt.Printf("Prov for %s\n", r.URL.Path[1:])
	fmt.Println(string(body))
	// do something with the POST data
	// likely convert to triples and write to some end point...

	w.WriteHeader(http.StatusNoContent)

	// w.Write([]byte("Thanks for your contribution"))  //  we are 204..  no need for body content
}
