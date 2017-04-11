package services

import (
	"bytes"
	// "fmt"
	"log"
	"time"

	sparql "github.com/knakk/sparql"
)

const queries = `
# Comments are ignored, except those tagging a query.

# The following gets the project data
# tag: CSDCO 
SELECT DISTINCT ?uri ?date ?lat ?long ?holeid
WHERE 
{ 
  ?uri rdf:type <http://opencoredata.org/id/voc/csdco/v1/CSDCOProject> . 
  ?uri <http://opencoredata.org/id/voc/csdco/v1/project> "{{.PROJID}}" . 
  ?uri <http://opencoredata.org/id/voc/csdco/v1/holeid> ?holeid .
  ?uri 	<http://opencoredata.org/id/voc/csdco/v1/date> ?date . 
  ?uri 	<http://www.w3.org/2003/01/geo/wgs84_pos#lat> ?lat .
  ?uri 	<http://www.w3.org/2003/01/geo/wgs84_pos#long> ?long .
}

# Get all the info on a HoleID from the CSDCO graph
# tag: CSDCOHoleID
SELECT *
WHERE 
{ 
  <{{.HOLEID}}>  ?p ?o .
}

# tag: my-query
SELECT *
WHERE {
  ?s ?p ?o
} LIMIT {{.Limit}} OFFSET {{.Offset}}

#tag: generic
SELECT *
FROM <http://data.oceandrilling.org/geolink/>
WHERE {
  <{{.URI}}> ?p ?o .
}

#tag: geolinkperson
prefix foaf: <http://xmlns.com/foaf/0.1/>
prefix owl: <http://www.w3.org/2002/07/owl#>
prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
prefix glview: <http://schema.geolink.org/1.0/base/main#> 
SELECT DISTINCT  ?name  ?dep
WHERE {
  VALUES ?target {<{{.URI}}>}
   ?target glview:hasFullName ?name .
   ?dep ?p ?target
}
GROUP BY ?name

#tag: blazetest
SELECT *
WHERE {
  <{{.URI}}> ?p ?o .
}

`

// CSDCOHoleIDInfo takes a project ID and returns the holeid URI's and lat long info in SPARQL results
// returns:  uri	date	lat	long	holeid
func CSDCOHoleIDInfo(holeid string) *sparql.Results {

	repo, err := sparql.NewRepo("http://localhost:9999/blazegraph/namespace/csdcov3/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("CSDCOHoleID", struct{ HOLEID string }{holeid})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res

}

// CSDCOProjectInfo takes a project ID and returns the holeid URI's and lat long info in SPARQL results
// returns:  uri	date	lat	long	holeid
func CSDCOProjectInfo(projid string) *sparql.Results {

	repo, err := sparql.NewRepo("http://localhost:9999/blazegraph/namespace/csdcov3/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("CSDCO", struct{ PROJID string }{projid})
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(res)
	return res

}

// GetGeoLinkResource for GeoLink All Hands Demo, remove afterwards, dont' want person specific version
func GetGeoLinkResource(uri string) *sparql.Results {
	repo, err := sparql.NewRepo("http://data.geolink.org/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("geolinkperson", struct{ URI string }{uri})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res //  Just a test return, would not be what I really return
}

// GetRDFResource takes a URI as an arugment and returns information about the RDF resource
func GetRDFResource(uri string) *sparql.Results {
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("generic", struct{ URI string }{uri})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// // Print loop testing
	// bindingsTest := res.Results.Bindings // map[string][]binding
	// fmt.Println("res.Resuolts.Bindings:")
	// for k, i := range bindingsTest {
	// 	fmt.Printf("At postion %v with %v and %v\n", k, i["pro"], i["vol"])
	// }

	// bindingsTest2 := res.Bindings() // map[string][]rdf.Term
	// fmt.Println("res.Bindings():")
	// for k, i := range bindingsTest2 {
	// 	fmt.Printf("At postion %v with %v \n", k, i)
	// }

	// solutionsTest := res.Solutions() // map[string][]rdf.Term
	// fmt.Println("res.Solutions():")
	// for k, i := range solutionsTest {
	// 	fmt.Printf("At postion %v with %v \n", k, i)
	// }

	return res //  Just a test return, would not be what I really return
}
