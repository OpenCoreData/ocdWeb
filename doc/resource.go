package doc

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"opencoredata.org/ocdWeb/services"
	"text/template"
)

// ResourceRender can pull data from graphdb and then display the results pubby like
// May start with virtuoso and migrate to cayley later...
func ResourceRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// log.Printf("for resource: %s\n", r.URL.Path)
	log.Printf("for resource: %s\n", vars["resourcepath"])

	res := services.GetRDFResource(fmt.Sprintf("http://opencoredata.org/id/resource/%s", vars["resourcepath"]))
	// fmt.Printf("%s", res)

	ht, err := template.New("some template").ParseFiles("templates/rdfResource.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	//  dataForTemplate := TemplateForDoc{Schema: result, CSVW: result2, Schemastring: string(jsonldtext), Csvwstring: string(csvwtext), UUID: vars["UUID"]}

	solutionsTest := res.Solutions() // map[string][]rdf.Term
	// make new map, pass to the template and call by key
	var solutionsMap map[string]string
	solutionsMap = make(map[string]string)
	for _, i := range solutionsTest {
		ps := fmt.Sprint(i["p"])
		os := fmt.Sprint(i["o"])
		solutionsMap[ps] = os
		fmt.Printf("KEY: %v \t\tVALUE %v \n", i["p"], i["o"])
	}

	err = ht.ExecuteTemplate(w, "T", solutionsMap) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}