package services

import (
	"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"bytes"
	"log"
	"net/http"
	"strings"
	// "encoding/json"
)

type MLCount struct {
	id      string `bson:"_id,omitempty"` // I don't really want the ID, so leave it lower case
	Measure string `json:"measure"`
	Leg     string `json:"leg"`
	Count   int    `json:"count"`
}

// Redirection handler
func GetGrid(w http.ResponseWriter, r *http.Request) {
	// call mongo and lookup the redirection to use...
	session, err := GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("aggregation_janusMLCountv2")

	var results []MLCount
	err = c.Find(nil).All(&results)
	if err != nil {
		log.Printf("Error calling aggregation_janusMLCount : %v", err)
	}

	// for _, elem := range results {
	// 	fmt.Printf("Leg: %s Measure: %s  Count %d\n", elem.Leg, elem.Measure, elem.Count)
	// }

	measures := UniqueMeasurements(results)
	fmt.Printf("Found %d items: %v\n", len(measures), measures)

	// this is stupid..  I should be able to make a struct and build this out...
	// but I don't see how the struct can't be anything other than something that
	// have every leg in it.
	// use a string buffer to build this out..
	var buffer bytes.Buffer
	buffer.WriteString("{\"rows\": [")
	for i, mes := range measures {
		// buffer.WriteString(fmt.Sprintf("{\"Measurement\": \"%s\",\n", mes)) // make a link to a collection too?
		buffer.WriteString(fmt.Sprintf("{\"Measurement\": \"<a href='/collections/measurement/%s'>%s</a>\",\n", mes, mes)) // make a link to a collection too?
		tosi := 1
		si := CountLegs(mes, results)
		for _, subelem := range results {
			if strings.Contains(subelem.Measure, mes) {
				// print the line
				buffer.WriteString(fmt.Sprintf("\"%s\" : \"<a href='/collections/%s/%s'>%d</a>\"", subelem.Leg, mes, subelem.Leg, subelem.Count))
				if tosi < si {
					buffer.WriteString(",\n")
				}
				tosi = tosi + 1
			}
		}
		buffer.WriteString("}")
		if i+1 < len(measures) {
			buffer.WriteString(",\n")
		}
	}
	buffer.WriteString("]}")

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%v", buffer.String())
}

func CountLegs(measure string, data []MLCount) int {
	var items []string
	for _, elem := range data {
		if strings.Contains(elem.Measure, measure) {
			items = append(items, elem.Measure)
		}
	}
	return len(items)
}

func UniqueMeasurements(data []MLCount) []string {
	var items []string
	for _, elem := range data {
		if !stringInSlice(elem.Measure, items) && elem.Measure != "" { // TODO  this null string this is about as stupid as the above
			items = append(items, elem.Measure)
		}
	}
	return items
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
