package do

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

# tag: projFDPs
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
     BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	  {
		?s <http://schema.org/about> ?ss .
		?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataSet> .
		?s <http://schema.org/distribution> ?dist .
		?dist <http://schema.org/contentUrl> ?res
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
	DOFDPs    []string
	DOResProj string
	Lat       string
	Long      string
	Bucket    string
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
	ct := r.Header.Get("Accept")
	vars := mux.Vars(r)
	oid := vars["ID"]
	var oi TypeCheck
	log.Printf("http://opencoredata.org/id/do/%s : %s \n", oid, ct)

	mc := utils.MinioConnectionDEV() // minio connection

	// Get a type check object from either a graph inquiry or an
	// object store inquiry.
	oi, err := getObjKern(oid) // returns TypeCheck{}
	if err != nil {
		log.Println(err)
		oi, err = getByID(oid, mc)
		if err != nil {
			log.Println(err)
		}
	}

	// Get the object but only if we didn't error above (no bucket)
	fo, err := mc.GetObject(oi.Bucket, oid, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// strings.Contains is kinda hackish, but why parse the ct into an
	// array and then check for a string when this should always work?
	if !strings.Contains(ct, "text/html") {
		fmt.Println("You don't seem to be a browse, good luck with this")

		// set default to octet stream?  but use stored if I have it
		w.Header().Set("Content-Type", "application/octet-stream")

		n, err := io.Copy(w, fo)
		if err != nil {
			log.Println("Issue with writing file to http response")
			log.Println(err)
		}
		log.Printf("Sent %d bytes\n", n)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// If we are headed to an HTML template, then read it to a buffer
		var buf bytes.Buffer
		nw := bufio.NewWriter(&buf)
		_, err = io.Copy(nw, fo)
		if err != nil {
			log.Println(err)
		}

		// if type == project or if type == do   (and so on)  should I do it this way?
		oi.DOResProj = buf.String() //  ?? what is this?  project JSON-LD object
		pf, _ := projResources(oid, "projfeatures")
		pd, _ := projResources(oid, "projdatasets")
		pp, _ := projResources(oid, "projFDPs")
		log.Printf("---- %v", pp)
		oi.DOFeature = pf
		oi.DOPkgMeta = pd
		oi.DOFDPs = pp

		// TODO ?  Should I make the template name associated with the bucketname?  Makes it easy to alter the templates
		ht, err := template.New("object template").ParseFiles("web/templates/objectDOResProj.html") // open and parse a template text file
		if err != nil {
			log.Printf("template parse failed: %s", err)
		}

		err = ht.ExecuteTemplate(w, "T", oi) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
		if err != nil {
			log.Printf("htemplate execution failed: %s", err)
		}
	}
}

func getObjKern(id string) (TypeCheck, error) {
	// Maps of types to buckets // the following needs to be in main and shared
	m := make(map[string]string)
	m["http://www.schema.org/DigitalDocument"] = "csdco-do-meta"
	m["http://opencoredata.org/voc/csdco/v1/Borehole"] = "csdco-do-resources"
	m["http://opencoredata.org/voc/csdco/v1/Project"] = "csdco-do-resources"
	m["http://schema.org/ResearchProject"] = "csdco-do-resources"
	m["http://schema.org/DataSet"] = "csdco-do-packages-meta"

	var results TypeCheck

	repo, err := services.BasementTS()
	if err != nil {
		log.Printf("%s\n", err)
		return results, err
	}

	// Set up the query collection
	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	// Do a type check call to see what the type of the object is
	q, err := bank.Prepare("typecheck", struct{ OID string }{id})
	if err != nil {
		log.Printf("%s\n", err)
	}

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	bindings := res.Results.Bindings // map[string][]rdf.Term

	if len(bindings) > 0 {
		results.Type = bindings[0]["type"].Value
		results.Graph = bindings[0]["graph"].Value
		results.Bucket = m[bindings[0]["type"].Value]
	} else {
		return results, errors.New("NoSuchGraphObject")
	}

	return results, err
}

func getByID(id string, mc *minio.Client) (TypeCheck, error) {
	// Maps of types to buckets // the following needs to be in main and shared
	m := make(map[string]string)
	// m["http://www.schema.org/DigitalDocument"] = "csdco-do-meta"
	// m["http://opencoredata.org/voc/csdco/v1/Borehole"] = "csdco-do-resources"
	// m["http://opencoredata.org/voc/csdco/v1/Project"] = "csdco-do-resources"
	// m["http://schema.org/ResearchProject"] = "csdco-do-resources"
	// m["http://schema.org/DataSet"] = "csdco-do-packages-meta"
	m["DigitalObject"] = "csdco-do"
	m["FrictionlessDataPackage"] = "csdco-do-packages"

	var results TypeCheck

	for i := range m {
		bucket := m[i]
		objectStat, e := mc.StatObject(bucket, id, minio.StatObjectOptions{})
		log.Println(m[i])
		if e != nil {
			continue
			// The following is nice code..  left for future reference
			// errResponse := minio.ToErrorResponse(e)
			// if errResponse.Code == "AccessDenied" {
			// 	return results, errors.New("AccessDenied")
			// }
			// if errResponse.Code == "NoSuchBucket" {
			// 	return results, errors.New("NoSuchBucket")
			// }
			// if errResponse.Code == "InvalidBucketName" {
			// 	return results, errors.New("InvalidBucketName")
			// }
			// if errResponse.Code == "NoSuchKey" {
			// 	return results, errors.New("NoSuchKey")
			// }
			// return results, errors.New("Unexpected Error") // or WTF...  either would work here
		}
		// if err is nil, it means I found it, so I could set and break!
		results.Bucket = m[i]
		results.Type = objectStat.ContentType
		log.Println(objectStat)
		break

	}

	// should do a check here to see if the above loop found something,
	// if not make a new error and send

	return results, nil
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
