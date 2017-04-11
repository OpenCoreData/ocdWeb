package services

import (
	"fmt"

	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"bytes"
	"log"
	"net/http"
	// "strings"
	// "encoding/json"
)

type CSDCO struct {
	LocationName           string
	LocationType           string
	Project                string
	LocationID             string
	Site                   string
	Hole                   string
	SiteHole               string
	OriginalID             string
	HoleID                 string
	Platform               string
	Date                   string
	WaterDepthM            string
	Country                string
	State_Province         string
	County_Region          string
	PI                     string
	Lat                    string
	Long                   string
	Elevation              string
	Position               string
	StorageLocationWorking string
	StorageLocationArchive string
	SampleType             string
	Comment                string
	mblfT                  string
	mblfB                  string
	MetadataSource         string
}

func CSDCOGrid(w http.ResponseWriter, r *http.Request) {

	// call mongo and lookup the redirection to use...
	session, err := GetMongoCon()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("csdco")

	var results []CSDCO
	err = c.Find(nil).All(&results)
	if err != nil {
		log.Printf("Error calling csdco : %v", err)
	}

	var buffer bytes.Buffer
	buffer.WriteString("{\"rows\": [")
	for index, subelem := range results {
		// add {
		buffer.WriteString("{")
		// put in each element of the struct
		buffer.WriteString(fmt.Sprintf("\"HoleIDF\" : \"<a target='_blank' href='/collections/%s/%s'>%s</a>\",", "csdco",
			subelem.HoleID, subelem.HoleID))
		buffer.WriteString(fmt.Sprintf("\"LocationName\" : \"%s\",", subelem.LocationName))
		buffer.WriteString(fmt.Sprintf("\"LocationType\" : \"%s\",", subelem.LocationType))
		buffer.WriteString(fmt.Sprintf("\"Project\" : \"<a target='_blank' href='/collections/%s/project/%s'>%s</a>\",", "csdco",
			subelem.Project, subelem.Project))
		buffer.WriteString(fmt.Sprintf("\"LocationID\" : \"%s\",", subelem.LocationID))
		buffer.WriteString(fmt.Sprintf("\"Site\" : \"%s\",", subelem.Site))
		buffer.WriteString(fmt.Sprintf("\"Hole\" : \"%s\",", subelem.Hole))
		buffer.WriteString(fmt.Sprintf("\"SiteHole\" : \"%s\",", subelem.SiteHole))
		buffer.WriteString(fmt.Sprintf("\"OriginalID\" : \"%s\",", subelem.OriginalID))
		buffer.WriteString(fmt.Sprintf("\"HoleID\" : \"%s\",", subelem.HoleID))
		buffer.WriteString(fmt.Sprintf("\"Platform\" : \"%s\",", subelem.Platform))
		buffer.WriteString(fmt.Sprintf("\"Date\" : \"%s\",", subelem.Date))
		buffer.WriteString(fmt.Sprintf("\"WaterDepthM\" : \"%s\",", subelem.WaterDepthM))
		buffer.WriteString(fmt.Sprintf("\"Country\" : \"%s\",", subelem.Country))
		buffer.WriteString(fmt.Sprintf("\"State_Province\" : \"%s\",", subelem.State_Province))
		buffer.WriteString(fmt.Sprintf("\"County_Region\" : \"%s\",", subelem.County_Region))
		buffer.WriteString(fmt.Sprintf("\"PI\" : \"%s\",", subelem.PI))
		buffer.WriteString(fmt.Sprintf("\"Lat\" : \"%s\",", subelem.Lat))
		buffer.WriteString(fmt.Sprintf("\"Long\" : \"%s\",", subelem.Long))
		buffer.WriteString(fmt.Sprintf("\"Elevation\" : \"%s\",", subelem.Elevation))
		buffer.WriteString(fmt.Sprintf("\"Position\" : \"%s\",", subelem.Position))
		buffer.WriteString(fmt.Sprintf("\"StorageLocationWorking\" : \"%s\",", subelem.StorageLocationWorking))
		buffer.WriteString(fmt.Sprintf("\"StorageLocationArchive\" : \"%s\",", subelem.StorageLocationArchive))
		buffer.WriteString(fmt.Sprintf("\"SampleType\" : \"%s\",", subelem.SampleType))
		buffer.WriteString(fmt.Sprintf("\"Comment\" : \"%s\",", subelem.Comment))
		buffer.WriteString(fmt.Sprintf("\"mblfT\" : \"%s\",", subelem.mblfT))
		buffer.WriteString(fmt.Sprintf("\"mblfB\" : \"%s\",", subelem.mblfB))
		buffer.WriteString(fmt.Sprintf("\"MetadataSource\" : \"%s\"", subelem.MetadataSource))
		buffer.WriteString("}")
		if index+1 < len(results) {
			buffer.WriteString(",\n")
		}
	}
	buffer.WriteString("]}")

	w.Header().Set("Content-type", "text/plain")
	fmt.Fprintf(w, "%v", buffer.String())
}
