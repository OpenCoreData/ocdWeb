package colls

import (
	"fmt"
	"log"
	"net/http"
	// "net/url"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdWeb/services"
)

type CSDCO struct {
	LocationName           string
	LocationType           string
	Project                string
	LocationID             string
	Site                   string
	Hole                   string
	SiteHole               string
	OriginalID             string
	HoleID                 string
	Platform               string
	Date                   string
	WaterDepthM            string
	Country                string
	State_Province         string
	County_Region          string
	PI                     string
	Lat                    string
	Long                   string
	Elevation              string
	Position               string
	StorageLocationWorking string
	StorageLocationArchive string
	SampleType             string
	Comment                string
	MblfT                  string
	MblfB                  string
	MetadataSource         string
}

// CSDCOOverview displays the overview matrix interface for the CSDCO holeids
func CSDCOOverview(w http.ResponseWriter, r *http.Request) {
	ht, err := template.New("some template").ParseFiles("templates/matrix_csdco_test.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", "results") //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

// CSDCOcollection shows the information page for a given project defined by the
// variable HoleID set in the request parameters.  It uses the
// template templates/catalog_csdco_new.html
func CSDCOcollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Connect to triplestore to get data via SPARQL bank call
	sparqlresults := services.CSDCOHoleIDInfo(vars["HoleID"])
	// var resultstest CSDCO

	// this is for the PROJ level, not the HOLE level.. move to another function
	uris := []string{}
	log.Println(sparqlresults.Results.Bindings)
	bindings := sparqlresults.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		log.Print(fmt.Sprintf("%v", i["uri"].Value))
		uris = append(uris, fmt.Sprintf("%v", i["uri"].Value))
	}

	log.Println(uris)

	// Connect to mongo and get the results
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("csdco")

	var results CSDCO
	err = c.Find(bson.M{"holeid": vars["HoleID"]}).One(&results)
	if err != nil {
		log.Printf("Error calling csdco : %v", err)
	}

	// log.Print(vars["HoleID"])
	// log.Print(results)

	ht, err := template.New("some template").ParseFiles("templates/catalog_csdco_new.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
