package do

import (
	"bytes"
	"log"

	"github.com/knakk/sparql"
	"opencoredata.org/ocdWeb/internal/services"
)

const uriquery = `
# Comments are ignored, except those tagging a query.

# tag: sparql
SELECT  ?o
WHERE {
  GRAPH ?graph
   {
    <{{.}}>  <http://www.w3.org/2000/01/rdf-schema#label>  ?o .
    }
}
`

func getURILabel(uri string) (string, error) {
	repo, err := services.BasementTS()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(uriquery)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("sparql", struct{ URI string }{uri})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	bindings := res.Results.Bindings // map[string][]rdf.Term
	// ??? 	b := res.Bindings()

	return bindings[0]["o"].Value, err
}
