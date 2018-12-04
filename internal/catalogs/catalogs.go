package catalogs

import (
	"fmt"
	"log"
	"net/http"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here

	"github.com/gorilla/mux"
	"opencoredata.org/ocdWeb/internal/utils"
)

// GeoLinkCatalog renders a catalog page
func GeoLinkCatalog(w http.ResponseWriter, r *http.Request) {
	log.Printf("For resource: %s\n", r.URL.String())

	// Go ahead and set up the template first..  if this fails we really just should get out nicely
	ht, err := template.New("geolink template").ParseFiles("templates/catalogGeoLink.html") //open and parse a template text file
	if err != nil {
		log.Printf("geolink template parse failed: %s", err)
	}

	// Read our VOiD document and get back an array of datasets
	dsa := utils.VoidReaderAll()

	// Set up some info in our catalog struct
	dc := utils.DataCatalog{ID: "http://opencoredata.org/catalogs/geolink", URL: "http://opencoredata.org/catalogs/geolink",
		Description: "A catalog of RDF graphs from Open Core Data for GeoLink that align to the GeoLink base ontology"}

	// Build a schema.org/DataCatalog entry
	cat, _ := utils.CatalogBuilder(dc, dsa)

	// create a function scoped struct to pass data to the template
	type ts struct {
		DSA []utils.VoidDataset
		SOD string
	}
	td := ts{DSA: dsa, SOD: string(cat)}

	err = ht.ExecuteTemplate(w, "T", td) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("Template execution failed: %s", err)
	}

}

// GeoLinkDataset renders a dataset page
func GeoLinkDataset(w http.ResponseWriter, r *http.Request) {
	log.Printf("For URL: %s\n", r.URL.String())

	vars := mux.Vars(r)
	// log.Printf("for resource: %s\n", r.URL.Path)
	log.Printf("For resource: %s\n", vars["resourcepath"])

	URI := fmt.Sprintf("http://opencoredata.org/catalog/geolink/dataset/%s", vars["resourcepath"])

	ht, err := template.New("geolink template").ParseFiles("templates/datasetGeoLink.html") //open and parse a template text file
	if err != nil {
		log.Printf("geolink template parse failed: %s", err)
	}

	dsa := utils.VoidReader(URI) // this gets all the datasets..  we want just one...

	cat, _ := utils.DsetBuilder(dsa[0])

	// create a function scoped struct to pass data to the template
	type ts struct {
		DSA utils.VoidDataset
		SOD string
	}
	td := ts{DSA: dsa[0], SOD: string(cat)}

	err = ht.ExecuteTemplate(w, "T", td) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("Template execution failed: %s", err)
	}

}
