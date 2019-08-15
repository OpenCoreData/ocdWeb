package do

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	// "net/url"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here

	"opencoredata.org/ocdWeb/internal/services"
	"opencoredata.org/ocdWeb/internal/utils"

	"github.com/gorilla/mux"
	"github.com/knakk/sparql"
	"github.com/minio/minio-go"
)

const queries = `
# Comments are ignored, except those tagging a query.

# tag: test
SELECT ?s ?name ?desc
WHERE {
  ?s ?p <http://opencoredata.org/id/do/cc7481953cacce428eda4f3ed11c96a4ea3b1114084acf29496c15908cb6dee4> .
  ?s <http://schema.org/name> ?name .
  ?s <http://schema.org/description> ?desc
}

# tag: typecheck
SELECT DISTINCT ?type ?graph
WHERE
{  GRAPH ?graph
    {
        <http://opencoredata.org/id/do/{{.OID}}> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
    }
}

# tag: projfeatures
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
    BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	 {
	    ?res <http://schema.org/about> ?ss .
	    ?res <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://opencoredata.org/voc/csdco/v1/Borehole> .
	 }
  }
}

# tag: projdatasets
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
     BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	  {
	     ?res <http://schema.org/about> ?ss .
	     ?res <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataSet> .
	  }
  }
}

# tag: getObject
SELECT  DISTINCT ?date ?mimetype ?type ?license ?filetype ?name ?desc ?related ?url ?text
WHERE {
  ?s ?p "{{.OID}}" .
  ?s2 ?p2 ?s .
  ?s2 <http://schema.org/dateCreated> ?date .
  ?s2 <http://schema.org/encodingFormat> ?mimetype .
  ?s2 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
  ?s2 <http://schema.org/license> ?license .
  ?s2 <http://schema.org/additionType> ?filetype .
  ?s2 <http://schema.org/name> ?name .
  ?s2 <http://schema.org/text> ?text .
  ?s2 <http://schema.org/description> ?desc .
  ?s2 <http://schema.org/isRelatedTo> ?related .
  ?s2 <http://schema.org/url> ?url
}
`

// TypeCheck is a list of parameters on a digital objects
type TypeCheck struct {
	Type      string
	Graph     string
	DOMeta    string
	DOPkgMeta []string
	DOFeature []string
	DOResProj string
	Lat       string
	Long      string
}

// ObjectKernel is a list of parameters on a digital objects
type ObjectKernel struct {
	Name     string
	Desc     string
	Date     string
	Mimetype string
	Type     string
	Licenses string
	Filetype string
	Related  string
	URL      string
	Text     string
}

// ObjectView collects an object and also does a SPARQL query for the type
// It uses the type to select a template and passes the package along to the
// template for web component rendering
func ObjectView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oid := vars["ID"]
	ct := r.Header.Get("Content-Type")

	log.Printf("Calling for http://opencoredata.org/id/do/%s\n", oid)
	log.Println(ct)

	repo, err := services.BasementTS()
	if err != nil {
		log.Printf("%s\n", err)
		return // TODO need to return to an ERROR page rather than error out, this is for triplestore down
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("typecheck", struct{ OID string }{oid})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	var results TypeCheck
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		results.Type = i["type"].Value
		results.Graph = i["graph"].Value
	}

	mc := utils.MinioConnectionDEV() // minio connection

	log.Println(results)

	// Maps of types to buckets // the following needs to be in main and shared
	m := make(map[string]string)
	m["http://www.schema.org/DigitalDocument"] = "csdco-do-meta"
	m["http://opencoredata.org/voc/csdco/v1/Borehole"] = "csdco-do-resources"
	m["http://opencoredata.org/voc/csdco/v1/Project"] = "csdco-do-resources"
	m["http://schema.org/ResearchProject"] = "csdco-do-resources"
	m["http://schema.org/DataSet"] = "csdco-do-packages-meta"

	log.Printf("%s:%s", m[results.Type], oid)

	// Get the object
	fo, err := mc.GetObject(m[results.Type], oid, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	if ct == "application/ld+json" {
		log.Println("write the object to w and return")
		n, err := io.Copy(w, fo)
		if err != nil {
			log.Println("Issue with writing file to http response")
			log.Println(err)
		}
		log.Printf("Copyed %d bytes\n", n)
		return
	}

	// Read the object into a buffer
	var buf bytes.Buffer
	nw := bufio.NewWriter(&buf)
	n, err := io.Copy(nw, fo)
	if err != nil {
		log.Println("Issue with writing file to buffer")
		log.Println(err)
	}
	log.Printf("Copyed %d bytes\n", n)

	// if type == project or if type == do   (and so on)  should I do it this way?
	results.DOResProj = buf.String()

	pf, _ := projResources(oid, "projfeatures")
	pd, _ := projResources(oid, "projdatasets")
	results.DOFeature = pf
	results.DOPkgMeta = pd

	// TODO ?  Should I make the template name associated with the bucketname?  Makes it easy to alter the templates
	ht, err := template.New("object template").ParseFiles("web/templates/objectDOResProj.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

func projResources(id, query string) ([]string, error) {
	repo, err := services.BasementTS()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare(query, struct{ ID string }{id})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	var r []string
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		r = append(r, i["res"].Value)
	}

	return r, err
}

// OLDObjectView looks for an object ID to search the graph on..  it
// then attempts to render a template with the results
func OLDObjectView(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oid := vars["ID"]

	repo, err := services.BasementTS()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("getObject", struct{ OID string }{oid})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// Connect to triplestore to get data via SPARQL bank call
	//sparqlresults := services.ObjectCall(vars["HoleID"])
	var results ObjectKernel

	// this is for the PROJ level, not the HOLE level.. move to another function
	// log.Println(res.Results.Bindings)
	bindings := res.Results.Bindings // map[string][]rdf.Term
	for _, i := range bindings {
		results.Name = i["name"].Value
		results.Desc = i["desc"].Value
		results.Date = i["date"].Value
		results.Mimetype = i["mimetype"].Value
		results.Type = i["type"].Value
		results.Licenses = i["license"].Value
		results.Filetype = i["filetype"].Value
		results.Related = i["related"].Value
		results.URL = i["url"].Value
		results.Text = i["text"].Value
	}

	ht, err := template.New("some template").ParseFiles("web/templates/object.html") // open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", results) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
