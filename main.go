package main

import (
	"code.google.com/p/gorilla/mux"
	"log"
	"net/http"
	"opencoredata.org/ocdWeb/colls"
	"opencoredata.org/ocdWeb/doc"
	"opencoredata.org/ocdWeb/dx"
	"opencoredata.org/ocdWeb/services"
	// _ "net/http/pprof"
)

func main() {
	// Common files, css, js, images, etc...
	rcommon := mux.NewRouter()
	rcommon.PathPrefix("/common/").Handler(http.StripPrefix("/common/", http.FileServer(http.Dir("./static"))))
	http.Handle("/common/", rcommon)

	// ParkingPage
	parking := mux.NewRouter()
	parking.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/ParkingPage"))))
	http.Handle("/", parking)

	// New Root to replace the old Root
	root := mux.NewRouter()
	root.PathPrefix("/root/").Handler(http.StripPrefix("/root/", http.FileServer(http.Dir("./static/Material"))))
	http.Handle("/root/", root)

	// Simpler services to support the web UI  (other services in ocdService)
	servroute := mux.NewRouter()
	servroute.HandleFunc("/services/grid", services.GetGrid)
	http.Handle("/services/", servroute)

	// recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/dataset/{UUID}", dx.Redirection)
	http.Handle("/id/", dxroute)

	//Browser by id redirection to doc  (gets a specific dataset)  http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B
	// http://opencoredata.org/doc/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	// http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B
	docroute := mux.NewRouter()
	docroute.HandleFunc("/doc/dataset/{UUID}", doc.UUIDRender)
	docroute.HandleFunc("/doc/dataset/{measurement}/{leg}/{site}/{hole}", doc.Render)
	http.Handle("/doc/", docroute)

	// Browse by collection   measurement leg site hole
	collections := mux.NewRouter()
	collections.HandleFunc("/collections/measurements/", colls.MLCounts)
	collections.HandleFunc("/collections/{measurements}/{leg}", colls.MLURLSets) //  called from the matrix page
	http.Handle("/collections/", collections)

	// Browse by expedition    leg site hole
	// expedition := mux.NewRouter()
	// expedition.PathPrefix("/expedition").Handler(http.StripPrefix("/expedition", http.FileServer(http.Dir("./static/ROOT"))))
	// http.Handle("/", expedition)

	// Later Browse options might include:  units, observations

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")

	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}

}
