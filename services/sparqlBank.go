package services

import (
	"bytes"
	"fmt"
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

#tag: paramcall
SELECT  ?uri ?name ?type ?column ?desc WHERE {  
	?uri <http://example.org/rdf/type> <http://opencoredata.org/id/voc/janus/v1/JanusQuerySet> .   
   ?uri     <http://opencoredata.org/id/voc/janus/v1/struct_name> "{{.MEASUREMENT}}" .  
   ?uri   <http://opencoredata.org/id/voc/janus/v1/go_struct_name> ?name .
   ?uri  <http://opencoredata.org/id/voc/janus/v1/go_struct_type> ?type .  
   ?uri    <http://opencoredata.org/id/voc/janus/v1/column_id> ?column  .
   ?uri    <http://opencoredata.org/id/voc/janus/v1/JanusMeasurement> ?jmes .  
   ?jmes  <http://opencoredata.org/id/voc/janus/v1/json_descript>  ?desc  
   }
   ORDER By (xsd:integer(?column))

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

#tag: allCSDCOProj
prefix csdco:  <http://opencoredata.org/id/voc/csdco/v1/> 
SELECT DISTINCT * 
WHERE { 
  ?s rdf:type csdco:CSDCOProject .
  ?s ?p ?o 
}

#tag: alljrsoexpeditions
prefix ocdjanus: <http://opencoredata.org/voc/janus/1/> 
prefix ocd: <http://opencoredata.org/voc/1/> 
SELECT DISTINCT  ?leg ?rvol ?label
WHERE {
  ?s  rdf:type  ocd:Drillsite .
  ?s  ocdjanus:leg ?leg .
  ?s  ocd:initialreportvolume ?rvol .
  
}
ORDER BY DESC(xsd:integer(?leg))

#tag: DEPRECATEDalljrsoexpeditions
prefix ocdjanus: <http://opencoredata.org/voc/janus/1/> 
prefix ocd: <http://opencoredata.org/voc/1/> 
SELECT DISTINCT ?s ?leg ?rvol ?cdata ?label
WHERE {
  ?s  rdf:type  ocd:Drillsite .
  ?s  ocdjanus:leg ?leg .
  ?s  ocd:initialreportvolume ?rvol .
  ?s  ocd:coredata ?cdata .
  ?s  rdfs:label ?label .
  
}
ORDER BY DESC(xsd:integer(?leg))

`

// connector function for the local sparql instance
func getLocalSPARQL() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://opencoredata.org/blazegraph/namespace/opencore/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

func getCSDCOSPARQL() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://opencoredata.org/blazegraph/namespace/opencore/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}

func getQuery(tag string) (string, error) {
	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)
	q, err := bank.Prepare(tag)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return q, err
}

// IN PROGRESS..   just a copy of JR version now
func AllCSDCOProjects() *sparql.Results {
	repo, err := getLocalSPARQL()
	q, err := getQuery("allCSDCOProj")

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res
}

// ParamCall  takes a janus measurement and pulls back the columns in it
func ParamCall(measurement string) *sparql.Results {
	repo, err := getLocalSPARQL()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("paramcall", struct{ MEASUREMENT string }{measurement})
	if err != nil {
		log.Printf("%s\n", err)
	}

	fmt.Printf("SPARQL: %s\n", q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res
}

// AllJRSOExpeditions returns all expeditons fro JRSO
// returns:  uri	date	lat	long	holeid
func AllJRSOExpeditions() *sparql.Results {
	// repo, err := sparql.NewRepo("http://opencore.dev/blazegraph/namespace/opencore/sparql",
	// 	sparql.Timeout(time.Millisecond*15000),
	// )
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// }

	// f := bytes.NewBufferString(queries)
	// bank := sparql.LoadBank(f)

	// q, err := bank.Prepare("alljrsoexpeditions")
	// if err != nil {
	// 	log.Printf("%s\n", err)
	// }

	repo, err := getLocalSPARQL()
	q, err := getQuery("alljrsoexpeditions")

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res
}

// CSDCOHoleIDInfo takes a project ID and returns the holeid URI's and lat long info in SPARQL results
// returns:  uri	date	lat	long	holeid
func CSDCOHoleIDInfo(holeid string) *sparql.Results {
	repo, err := getCSDCOSPARQL()
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
func CSDCOProjectInfo(projid string) (*sparql.Results, error) {
	repo, err := getCSDCOSPARQL()
	if err != nil {
		log.Printf("%s\n", err)
		return nil, err
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("CSDCO", struct{ PROJID string }{projid})
	if err != nil {
		log.Printf("%s\n", err)
		return nil, err
	}

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
		return nil, err
	}

	log.Println(res)
	return res, err

}

// GetGeoLinkResource data.geolink.org for GeoLink All Hands Demo, remove afterwards, dont' want person specific version
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

// GetRDFResource calls data.oceandrilling.org to takes a URI as an arugment and returns information about the RDF resource
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
