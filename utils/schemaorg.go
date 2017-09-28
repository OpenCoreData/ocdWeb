package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deiu/rdf2go"
	"github.com/kazarena/json-gold/ld"
)

// DataCatalog is a struct to hold metadata about data catalogs
type DataCatalog struct {
	ID          string
	URL         string
	Description string
}

// VoidDataset is a struct to hold items from a VOiD file that
// describe a dataset  https://developers.google.com/search/docs/data-types/datasets
type VoidDataset struct {
	ID                 string
	URL                string // type URL: Location of a page describing the dataset.
	Description        string // A short summary describing a dataset.
	Keywords           string // Keywords summarizing the dataset.
	Name               string // A descriptive name of a dataset (e.g., “Snow depth in Northern Hemisphere”)
	ContentURL         string
	AccrualPeriodicity string
	Issued             string
	License            string
	Publisher          string // Person, Org The name of the dataset creator (person or organization).
	Title              string
	DataDump           string
	Source             string
	LandingPage        string
	DownloadURL        string
	MediaType          string
	SameAs             string // type URL: Other URLs that can be used to access the dataset page.
	Version            string // The version number for this dataset.
	VariableMeasured   string // What does the dataset measure? (e.g., temperature, pressure)
	PublisherDesc      string
	PublisherName      string
	PublisherURL       string
	Latitude           string
	Longitude          string
}

// VoidReaderAll is the empty func 2 step for the case where I am not looking
// for a specific subject URI.  Could have used a string pointer to allow nil to be passed
// to the argument ref: https://stackoverflow.com/questions/32568977/golang-pass-nil-as-optional-argument-to-a-function
func VoidReaderAll() []VoidDataset {
	return VoidReader("")
}

// VoidReader takes a void document and extracts certain voc terms it's looking for
// into a struct.  It's a feeder for the DataSet and DataCatalog functions
func VoidReader(targetURI string) []VoidDataset {
	// Set a base URI  (is this the quad URI?)
	baseUri := "https://example.org/foo"

	// Create a new graph
	g := rdf2go.NewGraph(baseUri)
	file, _ := os.Open("./static/rdf/graph/void.ttl")
	nr := bufio.NewReader(file)

	// r is an io.Reader
	g.Parse(nr, "text/turtle")

	var vdsa []VoidDataset

	// Before the search see if we are looking for a specific subject URL
	var suri rdf2go.Term
	if targetURI == "" {
		suri = nil
	} else {
		suri = rdf2go.NewResource(targetURI)
	}

	triples := g.All(suri, nil, rdf2go.NewResource("http://rdfs.org/ns/void#Dataset"))
	for triple := range triples {
		var vds VoidDataset

		// fmt.Printf("Found the URI: %s \n", triples[triple].String())

		vds.ID = triples[triple].Subject.RawValue()  // hold what we are talking about
		vds.URL = triples[triple].Subject.RawValue() // The ID is the URL for this LOD case

		vds.Description = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/description"))
		// vds.Keywords = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/license"))
		vds.Name = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/title"))
		vds.ContentURL = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://www.w3.org/ns/dcat#downloadURL"))
		vds.AccrualPeriodicity = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/accrualPeriodicity"))
		vds.Issued = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/issued"))
		vds.License = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/license"))
		vds.Publisher = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/publisher"))
		vds.Title = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/title"))
		vds.DataDump = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://rdfs.org/ns/void#dataDump"))
		vds.Source = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://purl.org/dc/terms/source"))
		vds.LandingPage = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://www.w3.org/ns/dcat#landingPage"))
		vds.DownloadURL = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://www.w3.org/ns/dcat#downloadURL"))
		vds.MediaType = getObject(g, triples[triple].Subject, rdf2go.NewResource("http://www.w3.org/ns/dcat#mediaType"))
		vdsa = append(vdsa, vds)
	}

	return vdsa

}

func getObject(g *rdf2go.Graph, subjectURI, predicateURI rdf2go.Term) string {
	test := g.One(subjectURI, predicateURI, nil)

	if test != nil {
		return test.Object.RawValue()
	} else {
		return ""
	}
}

func DsetBuilder(dm VoidDataset) ([]byte, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	doc := map[string]interface{}{
		"@type": "Dataset",
		"@id":   dm.ID,
		"http://schema.org/url":              dm.URL,
		"http://schema.org/description":      dm.Description,
		"http://schema.org/keywords":         dm.Keywords,
		"http://schema.org/name":             dm.Name,
		"http://schema.org/variableMeasured": dm.VariableMeasured,
		"http://schema.org/distribution": map[string]interface{}{
			"@type": "DataDownload",
			"http://schema.org/contentUrl": dm.ContentURL,
		},
		"http://schema.org/publisher": map[string]interface{}{
			"@type": "Organization",
			"http://schema.org/description": dm.PublisherDesc,
			"http://schema.org/name":        dm.PublisherName,
			"http://schema.org/url":         dm.PublisherURL,
		},
		"http://schema.org/spatial": map[string]interface{}{
			"@type": "Place",
			"http://schema.org/geo": map[string]interface{}{
				"@type": "GeoCoordinates",
				"http://schema.org/latitude":  dm.Latitude,
				"http://schema.org/longitude": dm.Longitude,
			},
		},
	}

	context := map[string]interface{}{
		"@context": map[string]interface{}{
			"@vocab":  "http://schema.org/",
			"re3data": "http://example.org/re3data/0.1/",
		},
	}

	compactedDoc, err := proc.Compact(doc, context, options)
	if err != nil {
		fmt.Println("Error when compacting", err)
	}

	return json.MarshalIndent(compactedDoc, "", " ")
}

func CatalogBuilder(dc DataCatalog, dsa []VoidDataset) ([]byte, error) {

	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	// array of maps
	var dsArray []map[string]interface{}

	// we are basically going from stuct to map so we can pass the expected
	// type to the json-ld tools
	for _, v := range dsa {
		datasets := make(map[string]interface{})
		datasets["@type"] = "Dataset"
		datasets["description"] = v.Description
		datasets["url"] = v.URL
		dsArray = append(dsArray, datasets)
	}

	doc := map[string]interface{}{
		"@type": "DataCatalog",
		"@id":   dc.ID,
		"http://schema.org/url":         dc.URL,
		"http://schema.org/description": dc.Description,
		"http://schema.org/dataset":     dsArray,
	}

	context := map[string]interface{}{
		"@context": map[string]interface{}{
			"@vocab":  "http://schema.org/",
			"re3data": "http://example.org/re3data/0.1/",
		},
	}

	compactedDoc, err := proc.Compact(doc, context, options)
	if err != nil {
		fmt.Println("Error when compacting", err)
	}

	return json.MarshalIndent(compactedDoc, "", " ")
}
