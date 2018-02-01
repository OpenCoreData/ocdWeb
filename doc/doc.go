package doc

import (
	"log"
	"net/http"
	// "net/url"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"github.com/fatih/camelcase"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"opencoredata.org/ocdWeb/services"
	"opencoredata.org/ocdWeb/utils"
)

type Uriurl struct {
	Uri string
	Url string
}

// W3c csvw metadata structs
type CSVWMeta struct {
	Context      string       `json:"@context"`
	Dc_license   Dc_license   `json:"dc:license"`
	Dc_modified  Dc_modified  `json:"dc:modified"`
	Dc_publisher Dc_publisher `json:"dc:publisher"`
	Dc_title     string       `json:"dc:title"`
	Dcat_keyword []string     `json:"dcat:keyword"`
	TableSchema  TableSchema  `json:"tableSchema"`
	URL          string       `json:"url"`
}

type Dc_license struct {
	Id string `json:"@id"`
}

type Dc_modified struct {
	Type  string `json:"@type"`
	Value string `json:"@value"`
}

type Dc_publisher struct {
	Schema_name string     `json:"schema:name"`
	Schema_url  Schema_url `json:"schema:url"`
}

type Schema_url struct {
	Id string `json:"@id"`
}

type TableSchema struct {
	AboutURL   string    `json:"aboutUrl"`
	Columns    []Columns `json:"columns"`
	PrimaryKey string    `json:"primaryKey"`
}

type Columns struct {
	Datatype       string   `json:"datatype"`
	Dc_description string   `json:"dc:description"`
	Name           string   `json:"name"`
	Required       bool     `json:"required"`
	Titles         []string `json:"titles"`
}

// schema.org Dataset metadata structs            	Context             []interface{} `json:"@context"`
type SchemaOrgMetadata struct {
	Context             string          `json:"@context"`
	Type                string          `json:"@type"`
	Author              Author          `json:"author"`
	Description         string          `json:"description"`
	Distribution        Distribution    `json:"distribution"`
	GlviewDataset       string          `json:"glview:dataset"`
	GlviewKeywords      string          `json:"glview:keywords"`
	OpenCoreLeg         string          `json:"opencore:leg"`
	OpenCoreSite        string          `json:"opencore:site"`
	OpenCoreHole        string          `json:"opencore:hole"`
	OpenCoreMeasurement string          `json:"opencore:measurement"`
	Keywords            string          `json:"keywords"`
	Name                string          `json:"name"`
	Spatial             SpatialCoverage `json:"spatialCoverage"` // name miss-match due to it being held in that name in mongo
	URL                 string          `json:"url"`
}

type Author struct {
	Type        string `json:"@type"`
	Description string `json:"description"`
	Name        string `json:"name"`
	URL         string `json:"url"`
}

type Distribution struct {
	Type           string `json:"@type"`
	ContentURL     string `json:"contentUrl"`
	DatePublished  string `json:"datePublished"`
	EncodingFormat string `json:"encodingFormat"`
	InLanguage     string `json:"inLanguage"`
}

type SpatialCoverage struct {
	Type string `json:"@type"`
	Geo  Geo    `json:"geo"`
}

type Geo struct {
	Type      string `json:"@type"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type TemplateForDoc struct {
	Schema       SchemaOrgMetadata
	CSVW         CSVWMeta
	Schemastring string // template.JS
	Csvwstring   string // template.JS
	MeasureType  string
	UUID         string
	SPARQLTerm   string
}

// Render A document view function   Note NOT being used ...
// Called from main for measurement view  (need to FIX THIS)
// not sure if I want a M/L/S/H URL open or not at this time...
func Render(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	session, err := services.GetMongoCon()
	if err != nil {
		panic(err)
	}

	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("uniqueids")

	// Steps:   convert URL to URI and then go looking up the datasets

	// This is where I use the structs from ocdJanus
	URL := fmt.Sprintf("http://opencoredata.org/doc/dataset/%s/%s/%s", vars["leg"], vars["site"], vars["hole"])

	result := Uriurl{}
	err = c.Find(bson.M{"url": URL}).One(&result)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}

	log.Printf("doc:  %s", r.URL.Path)

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%s", result.Uri)
}

// UUIDRender displays a data set basedon the UID of the dataset
func UUIDRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// call mongo and lookup the redirection to use...
	session, err := services.GetMongoCon()
	if err != nil {
		// panic(err)
		log.Printf("MONGO:  error: %s", err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("schemaorg")
	c2 := session.DB("test").C("csvwmeta")

	// Get the schema.org data
	URI := fmt.Sprintf("http://opencoredata.org/id/dataset/%s", vars["UUID"])
	result := SchemaOrgMetadata{}
	err = c.Find(bson.M{"url": URI}).One(&result)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}

	// context setting hack
	result.Context = "http://schema.org"
	// jsonldtext, _ := json.MarshalIndent(result, "", " ") // results as embeddale JSON-LD

	// Get the CSVW  data to use in VariableMeasured
	result2 := CSVWMeta{}
	err = c2.Find(bson.M{"url": URI}).One(&result2)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}

	result2.Context = "http://www.w3.org/ns/csvw"
	csvwtext, _ := json.MarshalIndent(result2, "", " ") // TODO  REMOVE THIS     results as embeddale JSON-LD

	////////// start new jsonld code

	// Collect the parameters to add to the keywords....
	// IN LINE TEST to see if I can loop on CSVW parameters
	// could we add in a variable measured.
	var params []string
	for k, _ := range result2.TableSchema.Columns {
		// fmt.Println(result2.TableSchema.Columns[k].Name)
		params = append(params, result2.TableSchema.Columns[k].Name)
	}
	paramString := strings.Join(params[:], " ")

	dataSet := utils.VoidDataset{} // badly named struct...
	dataSet.ID = result.URL
	dataSet.URL = result.URL
	dataSet.Description = fmt.Sprintf("%s for ocean drilling expedition %s site %s hole %s", getJanusKeyword(result.Keywords), result.OpenCoreLeg, result.OpenCoreSite, result.OpenCoreHole)
	// dataSet.ContentURL = result.Distribution.ContentURL // like http://opencoredata.org/api/v1/documents/download/208_1262A_JanusThermalConductivity_VyaMsepM.csv
	dataSet.ContentURL = fmt.Sprintf("http://opencoredata.org/api/v1/documents/download/%s", result.Name)
	dataSet.Name = result.Name
	dataSet.Keywords = fmt.Sprintf("%s %s %s", paramString, getJanusKeyword(result.Keywords), result.Keywords)
	dataSet.PublisherName = result.Author.Name
	dataSet.PublisherURL = result.Author.URL
	dataSet.PublisherDesc = result.Author.Description
	dataSet.SameAs = result.URL
	dataSet.Latitude = result.Spatial.Geo.Latitude
	dataSet.Longitude = result.Spatial.Geo.Longitude
	dataSet.OpenCoreHole = result.OpenCoreHole
	dataSet.OpenCoreSite = result.OpenCoreSite
	dataSet.OpenCoreLeg = result.OpenCoreLeg

	// SPARQL version of variableMeasured section
	sparqlTerm := strings.Replace(strings.ToLower(getJanusKeyword(result.Keywords)), " ", "", -1)
	sparqlTerm = strings.Replace(sparqlTerm, "section", "", -1)
	sr := services.ParamCall(sparqlTerm)
	fmt.Printf("For query on %s \n %s \n", sparqlTerm, sr)

	// loop and create the variable measured entry  // NIGHTMARE to catch when I transition from mongo to RDF
	vma := []utils.VariableMeasured{}
	// for i := range result2.TableSchema.Columns {
	// 	vm := utils.VariableMeasured{}
	// 	vm.Description = result2.TableSchema.Columns[i].Dc_description
	// 	vm.UnitText = result2.TableSchema.Columns[i].Datatype
	// 	vm.Value = result2.TableSchema.Columns[i].Name
	// 	vma = append(vma, vm)
	// }

	for i := range sr.Results.Bindings {
		// fmt.Printf("Name value %s \n", sr.Results.Bindings[i]["name"].Value)
		// fmt.Printf("Name desc %s \n", sr.Results.Bindings[i]["desc"].Value)
		// fmt.Printf("Name type %s \n", sr.Results.Bindings[i]["type"].Value)
		// fmt.Printf("Name column %s \n", sr.Results.Bindings[i]["column"].Value)
		// fmt.Printf("Name uri %s \n", sr.Results.Bindings[i]["uri"].Value)
		vm := utils.VariableMeasured{}
		vm.Description = sr.Results.Bindings[i]["desc"].Value
		vm.UnitText = sr.Results.Bindings[i]["type"].Value
		vm.Value = sr.Results.Bindings[i]["name"].Value
		vm.URL = sr.Results.Bindings[i]["uri"].Value
		vma = append(vma, vm)
	}

	dataSet.VariableMeasured = vma

	newJsonLD, _ := utils.DsetBuilder(dataSet)
	// fmt.Printf("NEW\n %s\n\n", string(newJsonLD))

	////////// end new jsonld code

	// old schema.org print
	// fmt.Printf("OLD\n %s\n\n", string(jsonldtext))

	ht, err := template.New("some template").ParseFiles("templates/jrso_dataset_new.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	// need a simple function call to extract the "janus" keyword from the keyword string and toLower it and
	// pass it in this struct to use in the data types web component
	measureString := getJanusKeyword(result.Keywords)
	// dataForTemplate := TemplateForDoc{Schema: result, CSVW: result2, Schemastring: template.JS(string(jsonldtext)), Csvwstring: template.JS(string(csvwtext)), MeasureType: measureString, UUID: vars["UUID"]}
	dataForTemplate := TemplateForDoc{Schema: result, CSVW: result2, Schemastring: string(newJsonLD), Csvwstring: string(csvwtext), MeasureType: measureString, UUID: vars["UUID"], SPARQLTerm: sparqlTerm}

	err = ht.ExecuteTemplate(w, "T", dataForTemplate) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

// split the camel case string into various words
func getJanusKeyword(s string) string {
	ssplit := strings.Split(s, ",")
	var targetString string
	for _, element := range ssplit {
		if strings.Contains(strings.ToLower(element), "janus") {
			// targetString = strings.ToLower(element)

			splitted := camelcase.Split(element)          // hackish split of terms
			targetString = strings.Join(splitted[:], " ") //

		}
	}
	return strings.TrimSpace(targetString)
}

// TODO
// func getMeasurementInfo
