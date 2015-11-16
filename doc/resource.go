package doc

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// ResourceRender can pull data from Cayley and then display the results
// This will be my pubby....
func ResourceRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("for resource: %s\n", r.URL.Path)
	log.Printf("for resource: %s\n", vars["resourcepath"])
	fmt.Fprintf(w, "%s", "template for this resource coming")
}
