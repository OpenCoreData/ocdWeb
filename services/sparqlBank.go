package services

import (
	"bytes"
	// "fmt"
	sparql "github.com/knakk/sparql"
	"log"
	"time"
)

const queries = `
# Comments are ignored, except those tagging a query.

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

#tag: blazetest
SELECT *
WHERE {
  <{{.URI}}> ?p ?o .
}

`

// GetRDFResource takes a URI as an arugment and returns information about the RDF resource
func GetRDFResource(uri string) *sparql.Results {
	log.Printf("GetRDFResource: %s\n", uri)
	
	repo, err := sparql.NewRepo("http://data.oceandrilling.org/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	// q, err := bank.Prepare("my-query", struct{ Limit, Offset int }{10, 100})
	q, err := bank.Prepare("generic", struct{ URI string}{uri})
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
