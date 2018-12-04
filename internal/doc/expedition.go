package doc

import (
	"log"
	"net/http"

	// "net/url"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdWeb/internal/services"
)

type CruiseGL struct {
	Expedition    string `json:"expedition"`
	Cruisetype    string `json:"cruisetype"`
	Endportcall   string `json:"endportcall"`
	Operator      string `json:"operator"`
	Participant   string `json:"participant"`
	Program       string `json:"program"`
	Scheduler     string `json:"scheduler"`
	Startportcall string `json:"startportcall"`
	Legsitehole   string `json:"legsitehole"`
	Track         string `json:"track"`
	Vessel        string `json:"vessel"`
	Note          string `json:"note"` // use for rdfs:label
	Uri           string `json:"uri"`
	Rvol          string `json:"rvol"`
	Cdata         string `json:"cdata"`
}

type Feature struct {
	Uri                    string
	Lat                    string
	Long                   string
	Hole                   string
	Expedition             string
	Site                   string
	Program                string
	Waterdepth             string
	CoreCount              string
	Initialreportvolume    string
	Coredata               string
	Logdata                string
	Geom                   string
	Scientificprospectus   string
	CoreRecovery           string
	Penetration            string
	Scientificreportvolume string
	Expeditionsite         string
	Preliminaryreport      string
	CoreInterval           string
	PercentRecovery        string
	Drilled                string
	Vcdata                 string
	Note                   string
	Prcoeedingreport       string
}

type DataForTemplate struct {
	Schema       Feature
	SchemaString string
	Datasets     []SchemaOrgMetadata
}

func ShowFeature(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// call mongo and lookup the redirection to use...
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("features")

	// test for HOLE (not in all DSDP sites)
	var URI string
	_, ok := vars["HOLE"]
	if ok == true {
		URI = fmt.Sprintf("http://opencoredata.org/id/expedition/%s/%s/%s", vars["LEG"], vars["SITE"], vars["HOLE"])
	} else {
		URI = fmt.Sprintf("http://opencoredata.org/id/expedition/%s/%s", vars["LEG"], vars["SITE"])
	}

	result := Feature{}
	err = c.Find(bson.M{"uri": URI}).One(&result)
	if err != nil {
		log.Printf("ShowFeature URI lookup error: %v %v", err, vars)
	}
	jsonldtext, _ := json.MarshalIndent(result, "", " ") // results as embeddale JSON-LD

	ht, err := template.New("some template").ParseFiles("templates/feature_new.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	dataForTemplate := DataForTemplate{Schema: result, SchemaString: string(jsonldtext), Datasets: GetDatasets(vars["LEG"], vars["SITE"])}

	err = ht.ExecuteTemplate(w, "T", dataForTemplate)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}

//ShowExpedition is the handler for URL patterns: http://localhost/doc/expedition/28
func ShowExpedition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("expeditions")

	EXPEDITION := vars["LEG"]
	var results CruiseGL
	err = c.Find(bson.M{"expedition": EXPEDITION}).One(&results)
	if err != nil {
		log.Printf("Error calling for ShowExpeditions: %v", err)
	}

	// TODO  FIX!   The next 10 lines are hideous and should not need to be here.
	// I am stripping the trailing / that is in the graph data.
	// These slices are hideous..  some have trailing /'s that alter resolution
	lshSlice := strings.Split(results.Legsitehole, " ")

	// For each string in this slice need to check and remove any trailing /
	for k, _ := range lshSlice {
		if last := len(lshSlice[k]) - 1; last >= 0 && lshSlice[k][last] == '/' {
			lshSlice[k] = lshSlice[k][:last]
		}
	}

	// Make a new struct and put results and lshslice into it and pass it along
	type TemplateStruct struct {
		Cruise   CruiseGL
		LSHSlice []string
		Datasets []SchemaOrgMetadata
	}

	SendToTemplate := TemplateStruct{Cruise: results, LSHSlice: lshSlice, Datasets: GetDatasets(EXPEDITION, "")}

	ht, err := template.New("some template").ParseFiles("templates/expedition_new.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", SendToTemplate)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

// need to return an error too
func GetDatasets(Leg string, Site string) []SchemaOrgMetadata {
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("schemaorg")

	var results []SchemaOrgMetadata

	switch Site {
	case "":
		err = c.Find(bson.M{"opencoreleg": Leg}).All(&results)
	default:
		err = c.Find(bson.M{"opencoreleg": Leg, "opencoresite": Site}).All(&results)

	}

	if err != nil {
		log.Printf("Error calling for ShowExpeditions: %v", err)
		results = nil
	}

	return results

}

//todo this NOT what I want..   this is not Expeditions, this is features.....
func AllExpeditions(w http.ResponseWriter, r *http.Request) {

	sparqlresults := services.AllJRSOExpeditions()

	var results []CruiseGL

	// log.Println(sparqlresults.Results.Bindings)
	bindings := sparqlresults.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		var expedition CruiseGL // ?s ?leg ?rvol ?cdata ?label
		// log.Print(fmt.Sprintf("%v", i["s"].Value))
		expedition.Expedition = i["leg"].Value
		expedition.Rvol = i["rvol"].Value
		expedition.Note = i["label"].Value
		results = append(results, expedition)
	}

	ht, err := template.New("some template").ParseFiles("templates/expeditionsAll_new.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}

//todo this NOT what I want..   this is not Expeditions, this is features.....
func AllExpeditionsOLD(w http.ResponseWriter, r *http.Request) {
	// call mongo and lookup the redirection to use...
	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("expedire").C("expeditions")

	var results []CruiseGL
	// err = c.Find(bson.M{}).Sort("-expedition").All(&results)
	err = c.Find(bson.M{}).Sort("{ $natural: 1 }").All(&results)
	if err != nil {
		log.Printf("Error calling for AllExpeditions: %v", err)
	}

	ht, err := template.New("some template").ParseFiles("templates/expeditionsAll_new.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
