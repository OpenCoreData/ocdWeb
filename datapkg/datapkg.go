package datapkg

import (
	"archive/zip"
	"bytes"
	"encoding/json"
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
	SHA         string
	//FileList    []string
}

type DataPackage struct {
	Profile   string `json:"profile"`
	Resources []struct {
		Encoding string `json:"encoding"`
		Name     string `json:"name"`
		Path     string `json:"path"`
		Profile  string `json:"profile"`
	} `json:"resources"`
	SHA string
}

// ServePkg gets a pkg by hash and geenrates a landing page for it.
func ServePkg(w http.ResponseWriter, r *http.Request) {
	log.Println("In the get landing page for package function")

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

	b, err := getsdo(fo, oi.Size, ".", "metadata/schemaorg.json")
	if err != nil {
		fmt.Println(err)
	}

	// Get the package manafest
	dp, err := getsdo(fo, oi.Size, ".", "datapackage.json")
	if err != nil {
		fmt.Println(err)
		fmt.Println(dp)
	}

	// Go ahead and set up the template first..  if this fails we really just should get out nicely
	ht, err := template.New("data package template").ParseFiles("templates/grid_csdcoRes.html") //open and parse a template text file
	if err != nil {
		log.Printf("geolink template parse failed: %s", err)
	}

	td := parseSDO(string(b))
	td.SHA = vars["id"]
	dpp := parseManifest(dp)
	dpp.SHA = vars["id"]

	err = ht.ExecuteTemplate(w, "T", td) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	// TODO make this manefest section in the template
	err = ht.ExecuteTemplate(w, "M", dpp) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("Template execution failed: %s", err)
	}
}

// DownloadPkg for downloading the data package zip file, vs seeing the landing page
func DownloadPkg(w http.ResponseWriter, r *http.Request) {
	log.Println("In the get zip package function")
	mc := utils.MinioConnection() // minio connection
	vars := mux.Vars(r)
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
		log.Println(err)
	}
	log.Printf("Write file %s with bytes %d\n", f, n)
}

// DownloadPkgFile for downloading the data package zip file, vs seeing the landing page
func DownloadPkgFile(w http.ResponseWriter, r *http.Request) {
	log.Println("In the get resources from package function")

	mc := utils.MinioConnection() // minio connection

	vars := mux.Vars(r)
	f := fmt.Sprintf("%s.zip", vars["id"])

	// keys, ok := r.URL.Query()["key"]
	// if !ok || len(keys) < 1 {
	// 	log.Println("Url Param 'key' is missing")
	// 	return
	// }

	// key := keys[0]

	key := vars["resourcepath"]

	log.Println("Url Param 'key' is: " + string(key))
	log.Printf("Minio resource  is: %s\n", f)

	fo, err := mc.GetObject("packages", f, minio.GetObjectOptions{})
	if err != nil {
		log.Println("GetObject call error")
		fmt.Println(err)
	}

	oi, err := fo.Stat()
	if err != nil {
		log.Println("Stat call error")
		log.Println(err)
	}

	fmt.Println(oi)

	b, err := getsdo(fo, oi.Size, ".", string(key)) // TODO RENAME getsdo to getByteStream or something
	if err != nil {
		log.Println("getsdo call error")
		fmt.Println(err)
	}

	br := bytes.NewReader(b) // need []byte to be a reader....

	n, err := io.Copy(w, br)
	log.Printf("Write file %s with bytes %dd\n", f, n)
	if err != nil {
		log.Println("Issue with writing file to http response")
	}
}

func parseSDO(jld string) SDO {
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

func parseManifest(j []byte) DataPackage {
	dp := DataPackage{}

	if err := json.Unmarshal(j, &dp); err != nil {
		panic(err)
	}

	return dp
}

func getsdo(fo *minio.Object, offset int64, dest, target string) ([]byte, error) {
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
		if f.Name == target {
			b, err = extractToBytes(f)
			if err != nil {
				return nil, err
			}
		}
	}

	//TODO  while here, get the []string of files in the zip archive and return that too.

	return b, nil
}
