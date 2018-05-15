package datapkg

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alecthomas/template"
	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
	"github.com/tidwall/gjson"
	"opencoredata.org/ocdWeb/utils"
)

type SDO struct {
	FullJSON    string
	ID          string
	Description string
	ContentURL  string
	Keywords    string
	License     string
	Name        string
	PubName     string
	//FileList    []string
}

// ServePkg gets a pkg by hash and geenrates a landing page for it.
func ServePkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("Serving data package %s\n", vars["id"])

	mc := utils.MinioConnection() // minio connection

	f := fmt.Sprintf("%s.zip", vars["id"])

	fo, err := mc.GetObject("packages", f, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	oi, err := fo.Stat()
	if err != nil {
		log.Println("Issue with reading an object..  should I just fatal on this to make sure?")
	}

	fmt.Println(oi)

	b, err := getsdo(fo, oi.Size, ".")
	if err != nil {
		fmt.Println(err)
	}

	// Go ahead and set up the template first..  if this fails we really just should get out nicely
	ht, err := template.New("data package template").ParseFiles("templates/datapkg.html") //open and parse a template text file
	if err != nil {
		log.Printf("geolink template parse failed: %s", err)
	}

	td := parseSDO(string(b))

	err = ht.ExecuteTemplate(w, "T", td) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("Template execution failed: %s", err)
	}
}

// DownloadPkg
func DownloadPkg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	mc := utils.MinioConnection() // minio connection

	f := fmt.Sprintf("%s.zip", vars["id"])

	fo, err := mc.GetObject("packages", f, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
	}

	oi, err := fo.Stat()
	if err != nil {
		log.Println("Issue with reading an object..  should I just fatal on this to make sure?")
	}

	fmt.Println(oi)

	// ref http://www.mrwaggel.be/post/golang-transmit-files-over-a-nethttp-server-to-clients/
	// then I think the extrat to bytes from eblow will work and we can simply io.Copy the file trhough.

	//Send the headers
	// w.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	// w.Header().Set("Content-Type", FileContentType)
	// w.Header().Set("Content-Length", FileSize)

	n, err := io.Copy(w, fo)
	if err != nil {
		log.Println("Issue with writing file to http response")
	}
	log.Printf("Write file %s with bytes %d\n", f, n)

}

func parseSDO(jld string) SDO {
	log.Println("Parse elements from the schema.org")

	sdo := SDO{}

	sdo.FullJSON = jld
	sdo.ID = gjson.Get(jld, "@id").String()
	sdo.Description = gjson.Get(jld, "description").String()
	sdo.ContentURL = gjson.Get(jld, "distribution.contentUrl").String()
	sdo.Keywords = gjson.Get(jld, "keywords").String()
	sdo.License = gjson.Get(jld, "license").String()
	sdo.Name = gjson.Get(jld, "name").String()
	sdo.PubName = gjson.Get(jld, "publisher.name").String()
	// sd.FileList = []string passed to this function.

	return sdo
}

func getsdo(fo *minio.Object, offset int64, dest string) ([]byte, error) {
	r, err := zip.NewReader(fo, offset)
	if err != nil {
		return nil, err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractToBytes := func(f *zip.File) ([]byte, error) {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		b, err := ioutil.ReadAll(rc)
		if err != nil {
			return nil, err
		}

		return b, nil
	}

	var b []byte

	for _, f := range r.File {
		// fmt.Println(f.Name)
		if f.Name == "metadata/schemaorg.json" {
			b, err = extractToBytes(f)
			if err != nil {
				return nil, err
			}
		}
	}

	//TODO  while here, get the []string of files in the zip archive and return that too.

	return b, nil
}
