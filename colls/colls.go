package colls

import (
	"log"
	"net/http"
	// "net/url"
	// "fmt"
	"code.google.com/p/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

type MLSet struct {
	id      string   `bson:"_id,omitempty"` // I don't really want the ID, so leave it lower case
	Leg     string   `json:"leg"`
	Measure string   `json:"measure"`
	Urls    []string `json:"urls"`
}

func MLCounts(w http.ResponseWriter, r *http.Request) {
	ht, err := template.New("some template").ParseFiles("templates/collections.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", "results") //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

// needs to take a Leg and Measurement and return all data sets associated with it.
func MLURLSets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// call mongo and lookup the redirection to use...
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("aggregation_janusURLSet")

	log.Print(c.FullName)

	var results MLSet
	err = c.Find(bson.M{"measure": vars["measurements"], "leg": vars["leg"]}).One(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusURLSet : %v", err)
	}

	log.Print(results)

	ht, err := template.New("some template").ParseFiles("templates/measureSet.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}
