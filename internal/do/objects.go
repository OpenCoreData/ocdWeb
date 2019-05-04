package do

import (
	"bytes"
	"log"
	"net/http"

	// "net/url"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here

	"opencoredata.org/ocdWeb/internal/services"

	"github.com/gorilla/mux"
	"github.com/knakk/sparql"
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

// ObjectView looks for an object ID to search the graph on..  it
// then attempts to render a template with the results
func ObjectView(w http.ResponseWriter, r *http.Request) {
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
