package doc

import (
	"log"
	"net/http"
	// "net/url"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
)

type Uriurl struct {
	Uri string
	Url string
}

func Render(w http.ResponseWriter, r *http.Request) {

	// call mongo and lookup the redirection to use...
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("uniqueids")

	// Steps:   convert URL to URI and then go looking up the datasets

	// This is where I use the structs from ocdJanus
	URL := "http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B"
	result := Uriurl{}
	err = c.Find(bson.M{"url": URL}).One(&result)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}

	log.Printf("doc:  %s", r.URL.Path)

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%s", result.Uri)

}

func UUIDRender(w http.ResponseWriter, r *http.Request) {

	// call mongo and lookup the redirection to use...
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("uniqueids")

	// Steps:   convert URL to URI and then go looking up the datasets

	// This is where I use the structs from ocdJanus
	URI := "http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3"
	result := Uriurl{}
	err = c.Find(bson.M{"uri": URI}).One(&result)
	if err != nil {
		log.Printf("URL lookup error: %v", err)
	}

	log.Printf("doc:  %s", r.URL.Path)

	ht, err := template.New("some template").ParseFiles("templates/documentUUID.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", result) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
