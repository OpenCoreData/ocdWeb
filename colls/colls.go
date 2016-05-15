package colls

import (
	"log"
	"net/http"
	// "net/url"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdWeb/services"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here
)

type URLSet struct {
	id      string    `bson:"_id,omitempty"` // I don't really want the ID, so leave it lower case
	Leg     string    `json:"leg"`
	Measure string    `json:"measure"`
	Refdata []Refdata `json:"refdata"`
}

type Refdata struct {
	Url  string `json:"url"`
	Lat  string `json:"latitude"`
	Long string `json:"longitude"`
}

// schema.org Datacatalog struct
type SchemaDatacatalog struct {
	Context     string         `json:"@context"`
	Type        string         `json:"@type"`
	Author      SchemaAuthor   `json:"author"`
	Dataset     []ShemaDataset `json:"dataset"`
	Description string         `json:"description"`
	Name        string         `json:"name"`
	URL         string         `json:"url"`
}

type ShemaDataset struct {
	Type string `json:"@type"`
	URL  string `json:"url"`
}

type SchemaAuthor struct {
	Type        string `json:"@type"`
	Description string `json:"description"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type TemplateForColls struct {
	URLdata URLSet
	Schema  string
}

type TemplateForMeasurement struct {
	URLdata []URLSet
	Schema  string
}

// The template render doesn't do anything at time..  the .js in the page does all that for now
// Likely will do something wih the template later
func MLCounts(w http.ResponseWriter, r *http.Request) {
	ht, err := template.New("some template").ParseFiles("templates/matrix_jrso.html") //open and parse a template text file
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
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("aggregation_janusURLSet")

	var results URLSet
	err = c.Find(bson.M{"measure": vars["measurements"], "leg": vars["leg"]}).One(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusURLSet MLURLSet : %v", err)
	}

	// log.Print(results)
	// need to build simple metadata package around schema.org/DataCatalog
	authorInfo := SchemaAuthor{Type: "Organization", Name: "Joides Resolution Science Office",
		URL: "http://iodp.org", Description: "NSF funded operator for International Ocean Discvery Project"}
	dataSets := []ShemaDataset{}
	for _, d := range results.Refdata {
		dataSet := ShemaDataset{Type: "Dataset", URL: d.Url}
		dataSets = append(dataSets, dataSet)
	}
	dataCatalog := SchemaDatacatalog{Context: "http://schema.org",
		Type:        "DataCatalog",
		Author:      authorInfo,
		Dataset:     dataSets,
		Description: fmt.Sprintf("Data set for measurement %s and leg %s", vars["measurements"], vars["leg"]),
		Name:        fmt.Sprintf("%s%s", vars["measurements"], vars["leg"]),
		URL:         fmt.Sprintf("http://opencoredata.org/collections/measurement/%s/%s", vars["measurements"], vars["leg"])}

	schematext, _ := json.Marshal(dataCatalog) // .MarshalIndent(dataCatalog, "", " ")

	data := TemplateForColls{URLdata: results, Schema: string(schematext)}

	ht, err := template.New("some template").ParseFiles("templates/jrso_MS.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	// tmpl.Execute(out, template.HTML(`<b>World</b>`))

	err = ht.ExecuteTemplate(w, "T", data) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

func MesSets(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// call mongo and lookup the redirection to use...
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("aggregation_janusURLSet")

	var results []URLSet
	err = c.Find(bson.M{"measure": vars["measurements"]}).All(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusURLSet: in MesSet %v", err)
	}

	// log.Print(results)
	// need to build simple metadata package around schema.org/DataCatalog
	authorInfo := SchemaAuthor{Type: "Organization", Name: "Joides Resolution Science Office",
		URL: "http://iodp.org", Description: "NSF funded operator for International Ocean Discvery Project"}
	dataSets := []ShemaDataset{}
	for _, dp := range results {
		for _, d := range dp.Refdata {
			dataSet := ShemaDataset{Type: "Dataset", URL: d.Url}
			dataSets = append(dataSets, dataSet)
		}
	}
	dataCatalog := SchemaDatacatalog{Context: "http://schema.org",
		Type:        "DataCatalog",
		Author:      authorInfo,
		Dataset:     dataSets,
		Description: fmt.Sprintf("Data set for measurement %s ", vars["measurements"]),
		Name:        fmt.Sprintf("%s", vars["measurements"]),
		URL:         fmt.Sprintf("http://opencoredata.org/collections/measurement/%s", vars["measurements"])}

	schematext, _ := json.Marshal(dataCatalog) // .MarshalIndent(dataCatalog, "", " ")

	data := TemplateForMeasurement{URLdata: results, Schema: string(schematext)}

	ht, err := template.New("some template").ParseFiles("templates/jrso_ML.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	// tmpl.Execute(out, template.HTML(`<b>World</b>`))

	err = ht.ExecuteTemplate(w, "T", data) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

// The following function is deprectated.  I am leaving it in for a while in case I discovery some reference I need to deal with. 
//
// func LegSets(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)

// 	// call mongo and lookup the redirection to use...
// 	session, err := services.GetMongoCon()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()

// 	// Optional. Switch the session to a monotonic behavior.
// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("test").C("aggregation_janusURLSet")

// 	var results []URLSet
// 	err = c.Find(bson.M{"leg": vars["leg"]}).All(&results)
// 	if err != nil {
// 		log.Printf("Error calling aggregation_janusURLSet: in LegSet %v", err)
// 	}

// 	// log.Print(results)
// 	// need to build simple metadata package around schema.org/DataCatalog
// 	authorInfo := SchemaAuthor{Type: "Organization", Name: "Joides Resolution Science Office",
// 		URL: "http://iodp.org", Description: "NSF funded operator for International Ocean Discvery Project"}
// 	dataSets := []ShemaDataset{}
// 	for _, dp := range results {
// 		for _, d := range dp.Refdata {
// 			dataSet := ShemaDataset{Type: "Dataset", URL: d.Url}
// 			dataSets = append(dataSets, dataSet)
// 		}
// 	}
// 	dataCatalog := SchemaDatacatalog{Context: "http://schema.org",
// 		Type:        "DataCatalog",
// 		Author:      authorInfo,
// 		Dataset:     dataSets,
// 		Description: fmt.Sprintf("Data set for leg %s ", vars["leg"]),
// 		Name:        fmt.Sprintf("%s", vars["leg"]),
// 		URL:         fmt.Sprintf("http://opencoredata.org/doc/expedition/%s", vars["leg"])}

// 	schematext, _ := json.Marshal(dataCatalog) // .MarshalIndent(dataCatalog, "", " ")

// 	data := TemplateForMeasurement{URLdata: results, Schema: string(schematext)}

// 	ht, err := template.New("some template").ParseFiles("templates/jrso_ML.html") //open and parse a template text file
// 	if err != nil {
// 		log.Printf("template parse failed: %s", err)
// 	}

// 	// tmpl.Execute(out, template.HTML(`<b>World</b>`))

// 	err = ht.ExecuteTemplate(w, "T", data) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
// 	if err != nil {
// 		log.Printf("htemplate execution failed: %s", err)
// 	}
// }
