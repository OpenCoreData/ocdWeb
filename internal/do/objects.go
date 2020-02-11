package do

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	// "net/url"
	"text/template" // text not html since we don't want to escape our JSON-LD and we don't worry about the HTML autoescape here

	"opencoredata.org/ocdWeb/internal/services"

	"github.com/gorilla/mux"
	"github.com/knakk/sparql"
	"github.com/minio/minio-go"
)

// TypeCheck is a list of parameters on a digital objects
// This is what I send to the kernel
type TypeCheck struct {
	Type      string
	OID       string
	Graph     string
	DOMeta    string
	PkgsMeta  []string
	DOFeature []string
	DOFDPs    []string
	DOResProj string // JSON-LD ?
	Lat       string
	Long      string
	Bucket    string
	OrigName  string
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
func ObjectView(mc *minio.Client, w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Accept")
	vars := mux.Vars(r)
	oid := vars["ID"]
	var oi TypeCheck
	log.Printf("http://opencoredata.org/id/do/%s : %s \n", oid, ct)

	// Get a type check object from either a graph or object store inquiry.
	oi, err := getObjKern(oid) // returns TypeCheck{}
	if err != nil {
		log.Println(err)
		oi, err = getByID(oid, mc)
		if err != nil {
			log.Println(err)
		}
	}

	// If the object isn't know, we should error now..  no point.
	// Do a custom error and link to search?

	// Get the object but only if we didn't error above (no bucket)
	fo, err := mc.GetObject(oi.Bucket, oid, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// strings.Contains is kinda hackish, but why parse the ct into an
	// array and then check for a string when this should always work?
	// oi.Bucket == "csdco-do" is also hackish..  need to review this routing..
	// maybe the mux has some options I can leverage?
	log.Println(ct)

	// TODO  this whole routing issue is bad..   also..  my MUX likely can do this
	// for and SAVE ME FROM THIS!  (sub routes...)

	// TODO  (what about using accepts */*  ?)
	if strings.Contains(ct, "text/html") && oi.Bucket == "csdco-do-packages" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// If we are headed to an HTML template, then read be object to a buffer
		// var buf bytes.Buffer
		// nw := bufio.NewWriter(&buf)
		// _, err = io.Copy(nw, fo)
		// if err != nil {
		// 	log.Println(err)
		// }

		log.Printf("Bucket: %s , Type: %s \n", oi.Bucket, oi.Type)
		t := "web/templates/objectDODataSet.html"

		// This is wrong.. needs to be the schema.org data graph for this
		// resource..  not the resource itself
		//  ?? what is this?  project JSON-LD object
		pg, _ := projResources(oid, "sdograph") // get the meta data graph for a DO OID
		u, err := url.Parse(pg[0])
		if err != nil {
			log.Fatal(err)
		}

		pa := strings.Split(u.Path, "/")

		log.Println("------------ sdo object I need -----------------")
		fmt.Println(pa[len(pa)-1]) // need to pull this object
		// take my OBJECT URI, parse the OID, get the file..  locate to oi.DOResProj (make this variable name more generic)

		// Get a type check object from either a graph or object store inquiry.
		oi, err := getObjKern(pa[len(pa)-1]) // returns TypeCheck{}
		if err != nil {
			log.Println(err)
			oi, err = getByID(pa[len(pa)-1], mc)
			if err != nil {
				log.Println(err)
			}
		}

		log.Println(oi.Bucket)
		fo, err := mc.GetObject(oi.Bucket, pa[len(pa)-1], minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err)
		}

		log.Println("------------ object (contents) -----------------")
		var b bytes.Buffer
		bw := bufio.NewWriter(&b)

		_, err = io.Copy(bw, fo)
		if err != nil {
			log.Println(err)
		}

		log.Println(string(b.Bytes()))

		oi.DOResProj = string(b.Bytes())

		// TODO REMOVE..  just pass this URL as []string
		pp, _ := projResources(oid, "projFDPs")
		pp = append(pp, fmt.Sprintf("http://opencoredata.org/id/do/%s", oid))

		// todo deal with these errors!
		log.Printf("Project FDPs:  %v", pp)
		// oi.DOFeature = pf
		// oi.PkgsMeta = pd
		oi.DOFDPs = pp
		oi.OID = oid

		// TODO ?  Should I make the template name associated with the bucketname?  Makes it easy to alter the templates
		ht, err := template.New("object template").ParseFiles(t) // open and parse a template text file
		if err != nil {
			log.Printf("template parse failed: %s", err)
		}

		err = ht.ExecuteTemplate(w, "T", oi) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
		if err != nil {
			log.Printf("htemplate execution failed: %s", err)
		}

		// send the bytes
		// n, err := io.Copy(w, fo)
		// if err != nil {
		// 	log.Println("Issue with writing file to http response")
		// 	log.Println(err)
		// }
		// log.Printf("NEW :   Sent %d bytes\n", n)
	} else if !strings.Contains(ct, "text/html") || oi.Bucket == "csdco-do" || oi.Bucket == "csdco-do-packages" {
		fmt.Println("You don't seem to be a browse, good luck with this")

		// set default to octet stream?  but use stored if I have it
		if oi.Type != "" {
			w.Header().Set("Content-Type", oi.Type)
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", oi.OrigName))
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		// unless we see a particular request type for a particular bucket just copy
		// and send.   How do we deal with the I want ZIP?   How do we register
		// what an objct can offer?  Just more content types?
		// Then go services are based on content type and service offering?

		// Service routers
		// if ct == zip and  oi.Type is frictionless data packages
		// bytes = makeFDPZIP(FDP data package json)

		// if ct == geojson and oi.Type(?) research project
		// make geojson

		// send the bytes
		n, err := io.Copy(w, fo)
		if err != nil {
			log.Println("Issue with writing file to http response")
			log.Println(err)
		}

		log.Printf("Sent %d bytes\n", n)
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// If we are headed to an HTML template, then read be object to a buffer
		var buf bytes.Buffer
		nw := bufio.NewWriter(&buf)
		_, err = io.Copy(nw, fo)
		if err != nil {
			log.Println(err)
		}

		log.Printf("Bucket: %s , Type: %s \n", oi.Bucket, oi.Type)
		t := "web/templates/objectDOResProj.html"

		// TODO these are only for research project type ????
		// if type == project or if type == do   (and so on)  should I do it this way?
		oi.DOResProj = buf.String() //  ?? what is this?  project JSON-LD object
		pf, _ := projResources(oid, "projfeatures")
		pd, _ := projResources(oid, "projdatasets")
		pp, _ := projResources(oid, "projFDPs") // need to be FDP directly
		// todo deal with these errors!
		log.Printf("Project FDPs:  %v", pp)
		oi.DOFeature = pf
		oi.PkgsMeta = pd
		oi.DOFDPs = pp
		oi.OID = oid

		// TODO ?  Should I make the template name associated with the bucketname?  Makes it easy to alter the templates
		ht, err := template.New("object template").ParseFiles(t) // open and parse a template text file
		if err != nil {
			log.Printf("template parse failed: %s", err)
		}

		err = ht.ExecuteTemplate(w, "T", oi) // substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
		if err != nil {
			log.Printf("htemplate execution failed: %s", err)
		}
	}
}

// // IsAType template function to test its type?
// func (d TypeCheck) IsAType(restype string) bool {
// 	if restype == "http://schema.org/ResearchProject" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

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
		results.OrigName = objectStat.Metadata["X-Amz-Meta-Filename"][0]
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

	log.Printf("\n %s \n", q)

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
