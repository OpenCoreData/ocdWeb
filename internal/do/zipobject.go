package do

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go"
)

// FDP is the Frictionless Data Package stuct.
type FDP struct {
	Name        string         `json:"name"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Sources     []FDPSources   `json:"sources"`
	Resources   []FDPResources `json:"resources"`
}

// FDPSources is the collection of sources for a FDP package
type FDPSources struct {
	Name string `json:"name"`
	Web  string `json:"web"`
}

// FDPResources are the data resources that are part of the FDP instance
// ?description ?name ?license ?encodingFormat ?url ?type ?additionType ?dateCreated ?identifier
type FDPResources struct {
	URL            string `json:"path"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	License        string `json:"licenses"`
	EncodingFormat string `json:"mediatype"`
	Type           string `json:"type"`
	AdditionalType string `json:"additionalType"`
	// DateCreated    string `json:""`
	// Identifuer     string `json:""`
}

// DownloadPkg https://stackoverflow.com/questions/46791169/create-serve-over-http-a-zip-file-without-writing-to-disk
func DownloadPkg(mc *minio.Client, w http.ResponseWriter, r *http.Request) {
	log.Println("In the get zip package function")
	//ct := r.Header.Get("Accept")
	vars := mux.Vars(r)
	oid := vars["ID"]

	//Send the headers
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", oid))
	// w.Header().Set("Content-Length", FileSize)

	// Make a new writer
	zw := zip.NewWriter(w)

	// TODO parse the JSON and get all the object
	// TODO rewrite the JSON to local form, not remove form

	// Get our package object
	// Minio object
	fo, err := mc.GetObject("csdco-do-packages", oid, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	fdp := FDP{}
	var b bytes.Buffer
	bw := bufio.NewWriter(&b)

	_, err = io.Copy(bw, fo)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(b.Bytes(), &fdp)
	if err != nil {
		log.Println(err)
	}

	for i := range fdp.Resources {
		o, err := urlPath(fdp.Resources[i].URL)
		if err != nil {
			log.Println(err)
		}
		log.Println(o)
		// range on the stuct items that are objects and pass to addObject()
		// sub in o
		// TODO  WARNING on this line..   a place holder till deployed to where the DO byte stream are
		err = addObject(mc, "000003a5ee30630237ae9690fd10a576b8bfb3d6c3e2ce541924522ef5b69f2c", fdp.Resources[i].Name, zw)
		if err != nil {
			log.Println(err)
		}
	}

	// Rewrite datapackage.json to relative links
	rel, err := fdpURLToRel(fdp)
	if err != nil {
		log.Println(err)
	}
	// load it to the zip
	cf, err := zw.Create("datapackage.json")
	if err != nil {
		log.Fatal(err)
	}

	// copy the object contents to the zip Writer
	n, err := io.Copy(cf, strings.NewReader(rel))
	log.Printf("Copied %d bytes\n", n)
	if err != nil {
		log.Fatal(err)
	}

	// close the zip Writer to flush the contents to the ResponseWriter
	err = zw.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// caution..  never pass by reference here ... might not be good
func fdpURLToRel(f FDP) (string, error) {
	for i := range f.Resources {
		f.Resources[i].URL = fmt.Sprintf("./data/%s", f.Resources[i].Name)
	}

	j, err := json.MarshalIndent(f, "", " ")
	if err != nil {
		log.Println(err)
	}

	return string(j), nil
}

func urlPath(s string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	p := u.Path
	sp := strings.Split(p, "/")
	oid := sp[len(sp)-1]

	return oid, err
}

func addObject(mc *minio.Client, oid, name string, zw *zip.Writer) error {
	// Minio object
	fo, err := mc.GetObject("csdco-do", oid, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	// write straight to the http.ResponseWriter so can avoid local marshalling
	cf, err := zw.Create(fmt.Sprintf("data/%s", name))
	if err != nil {
		log.Fatal(err)
	}

	// copy the object contents to the zip Writer
	n, err := io.Copy(cf, fo)
	log.Printf("Copied %d bytes\n", n)
	if err != nil {
		log.Fatal(err)
	}

	return err
}
