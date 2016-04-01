package doc

import (
	"log"
	"net/http"
	// "net/url"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"opencoredata.org/ocdWeb/services"
	"strings"
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
	Note          string `json:"note"`
	Uri           string `json:"uri"`
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

	ht, err := template.New("some template").ParseFiles("templates/feature.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	dataForTemplate := DataForTemplate{Schema: result, SchemaString: string(jsonldtext), Datasets: GetDatasets(vars["LEG"], vars["SITE"])}

	err = ht.ExecuteTemplate(w, "T", dataForTemplate)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}

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

	// These slices are hideous..  some have trailing /'s that alter resolution
	lshSlice := strings.Split(results.Legsitehole, " ")

	// Make a new struct and put results and lshslice into it and pass it along
	type TemplateStruct struct {
		Cruise   CruiseGL
		LSHSlice []string
		Datasets []SchemaOrgMetadata
	}

	SendToTemplate := TemplateStruct{Cruise: results, LSHSlice: lshSlice, Datasets: GetDatasets(EXPEDITION, "")}

	ht, err := template.New("some template").ParseFiles("templates/expedition.html")
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
	err = c.Find(bson.M{}).All(&results)
	if err != nil {
		log.Printf("Error calling for AllExpeditions: %v", err)
	}

	ht, err := template.New("some template").ParseFiles("templates/expeditionsAll.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
