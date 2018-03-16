package colls

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
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
	URI                    string
}

type CSDCOResultSet struct {
	Project string
	CSDCO   []CSDCO
}

type Abstract struct {
	ID          string               `bson:"id,omitempty"` // this is the ID, not the mongo _id
	Title       string               `bson:"title,omitempty"`
	Abstract    string               `bson:"abstract,omitempty"`
	Tags        []string             `bson:"tags,omitempty"`
	Identifiers AbstractsIdentifiers `bson:"identifiers,omitempty"`
}

type AbstractsIdentifiers struct {
	Issn string `bson:"issn,omitempty"`
	Doi  string `bson:"doi,omitempty"`
	Isbn string `bson:"isbn,omitempty"`
	Pmid string `bson:"pmid,omitempty"`
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

// CSDCOAbstract get the abstract info for a project
func CSDCOAbstract(w http.ResponseWriter, r *http.Request) {

	log.Println("CSDCO Abstract handler")
	vars := mux.Vars(r)

	// call mongo and lookup the redirection to use...
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("abstracts").C("csdco")

	var results Abstract
	err = c.Find(bson.M{"id": vars["ID"]}).One(&results)
	if err != nil {
		log.Printf("Error calling CSDCO abstract mongo : %v", err)
	}

	ht, err := template.New("abstract template").Funcs(template.FuncMap{"hasField": hasField}).ParseFiles("templates/csdco_abstract.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}

func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// CSDCOProjectInfo provides information via SPARQL on CSDCO Projects
func CSDCOProjectInfo(w http.ResponseWriter, r *http.Request) {

	// Connect to triplestore to get data via SPARQL bank call
	vars := mux.Vars(r)
	sparqlresults, _ := services.CSDCOProjectInfo(vars["ProjectID"])

	// this is for the PROJ level, not the HOLE level.. move to another function
	//uris := []string{}
	var results []CSDCO
	resultset := CSDCOResultSet{}

	if sparqlresults != nil {
		bindings := sparqlresults.Results.Bindings // map[string][]rdf.Term
		for _, i := range bindings {
			var result CSDCO
			// log.Print(fmt.Sprintf("%v", i["uri"].Value))
			result.HoleID = i["holeid"].Value
			result.Lat = i["lat"].Value
			result.Long = i["long"].Value
			result.Date = i["date"].Value
			result.URI = i["uri"].Value
			results = append(results, result) // fmt.Sprintf("%v", i["uri"].Value))
		}

		//log.Println(results)
		resultset.Project = vars["ProjectID"]
		resultset.CSDCO = results
	}

	ht, err := template.New("some template").ParseFiles("templates/csdco_projects.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", resultset) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
