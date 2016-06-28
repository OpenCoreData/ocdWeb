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

    URI := fmt.Sprintf("http://opencoredata.org/id/resource/%s", vars["resourcepath"])

	res := services.GetRDFResource(URI)
	// fmt.Printf("%s", res)

	ht, err := template.New("some template").ParseFiles("templates/rdfResource_new.html") //open and parse a template text file
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

	solutionsMap["URI"] = URI

	// add the resource ID to the Map too

	err = ht.ExecuteTemplate(w, "T", solutionsMap) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}
}

 // for GeoLink All Hands Demo, remove afterwards, dont' want person specific version
 func PersonResourceRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// log.Printf("for resource: %s\n", r.URL.Path)
	log.Printf("for person resource: %s\n", vars["resourcepath"])

    URI := fmt.Sprintf("http://opencoredata.org/id/resource/people/%s", vars["resourcepath"])

	res := services.GetGeoLinkResource(URI)
	// fmt.Printf("%s", res)

	ht, err := template.New("some template").ParseFiles("templates/rdfPersonResource_new.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	//  dataForTemplate := TemplateForDoc{Schema: result, CSVW: result2, Schemastring: string(jsonldtext), Csvwstring: string(csvwtext), UUID: vars["UUID"]}

	solutionsTest := res.Solutions() // map[string][]rdf.Term
	// make new map, pass to the template and call by key
	var solutionsMap map[string]string
	solutionsMap = make(map[string]string)
	for _, i := range solutionsTest {
		ps := fmt.Sprint(i["name"])
		os := fmt.Sprint(i["dep"])
		solutionsMap[os] = ps
		fmt.Printf("KEY: %v \t\tVALUE %v \n", i["name"], i["dep"])
	}

    var infoMap map[string]string
	infoMap = make(map[string]string)
	infoMap["URI"] = URI

	// add the resource ID to the Map too
	// this is ugly..  it all has to go.. this is an EarthCube demo hack to deal with later

	err = ht.ExecuteTemplate(w, "A", infoMap) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", solutionsMap) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}


err = ht.ExecuteTemplate(w, "Z", infoMap) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}