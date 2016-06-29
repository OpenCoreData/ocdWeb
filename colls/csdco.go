package colls

import (
	"log"
	"net/http"
	// "net/url"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdWeb/services"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here
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

func CSDCOcollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

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

	log.Print(vars["HoleID"])
	log.Print(results)

	ht, err := template.New("some template").ParseFiles("templates/catalog_csdco_new.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
