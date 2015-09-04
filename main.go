package main

/**
Web Resident Application for Image Strips

Need to take a URI like: http://data.oceandrilling.org/imagestrips/stripset/_05xx/0504B/0504B256M_1
and display it.
go get code.google.com/p/gorilla/mux
go get code.google.com/p/graphics-go/graphics
go get code.google.com/p/go.imagesdss
**/

import (
	"code.google.com/p/gorilla/mux"
	"log"
	"net/http"
	// _ "net/http/pprof"
)

func main() {
	// Common files, css, js, images, etc...
	rcommon := mux.NewRouter()
	rcommon.PathPrefix("/common/").Handler(http.StripPrefix("/common", http.FileServer(http.Dir("./static"))))
	http.Handle("/common/", rcommon)

	// ROOT
	root := mux.NewRouter()
	//root.Headers("Access-Control-Allow-Origin", "*")
	//root.Headers("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	root.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/ROOT"))))
	http.Handle("/", root)

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")

	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}

}
