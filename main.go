package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"opencoredata.org/ocdWeb/colls"
	"opencoredata.org/ocdWeb/doc"
	"opencoredata.org/ocdWeb/dx"
	"opencoredata.org/ocdWeb/services"
	"opencoredata.org/ocdWeb/voc"

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
	servroute.HandleFunc("/services/csdco", services.CSDCOGrid)
	http.Handle("/services/", servroute)

	// Recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// Example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/dataset/{UUID}", dx.Redirection)
	dxroute.HandleFunc(`/id/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, dx.RDFRedirection)
	http.Handle("/id/", dxroute)

	// Deal with void...  (show void..  allow .rdf file downloads)
	rdfdocs := mux.NewRouter()
	rdfdocs.PathPrefix("/rdf/").Handler(http.StripPrefix("/rdf/", http.FileServer(http.Dir("./static/rdf"))))
	http.Handle("/rdf/", rdfdocs)

	// Display Vocabulary entries.  A simple human view..
	// For machines, check for accepts headers?
	// no 303 for these?
	vocroute := mux.NewRouter()
	vocroute.PathPrefix("/voc/1/ocdSKOS.ttl").Handler(http.StripPrefix("/voc/", http.FileServer(http.Dir("./static/voc"))))
	vocroute.PathPrefix("/voc/janus/1/ocdJanusSKOS.ttl").Handler(http.StripPrefix("/voc/janus", http.FileServer(http.Dir("./static/voc/janus"))))
	vocroute.HandleFunc("/voc/janus/{version}/{term}", voc.VocJanus)
	vocroute.HandleFunc("/voc/{version}/{term}", voc.VocCore)
	http.Handle("/voc/", vocroute)

	//Browser by id redirection to doc  (gets a specific dataset)  http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B
	docroute := mux.NewRouter()
	docroute.HandleFunc("/doc/dataset/{UUID}", doc.UUIDRender)
	docroute.HandleFunc(`/doc/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, doc.ResourceRender)
	docroute.HandleFunc("/doc/dataset/{measurement}/{leg}/{site}/{hole}", doc.Render)
	http.Handle("/doc/", docroute)

	// Browse by collection   measurement leg site hole
	// Later Browse options might include:  units, observations. geologic time
	// TODO  worry about namespace collision here...  (need operator ID ?)
	collections := mux.NewRouter()
	// collections.HandleFunc("/collections", colls.Landing)
	collections.HandleFunc("/collections/matrix", colls.MLCounts)
	// collections.HandleFunc("/collections/januslegs", colls.JanusLegs)
	// collections.HandleFunc("/collections/janusmeasurements", colls.JanusMeasurements)
	collections.HandleFunc("/collections/csdco", colls.CSDCOOverview)
	collections.HandleFunc("/collections/csdco/{HoleID}", colls.CSDCOcollection) //  landing page for collection of files with a HoleID
	collections.HandleFunc("/collections/measurement/{measurements}", colls.MesSets)
	collections.HandleFunc("/collections/leg/{leg}", colls.LegSets)
	collections.HandleFunc("/collections/{measurements}/{leg}", colls.MLURLSets) //  called from the jrso matrix page
	http.Handle("/collections/", collections)

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")

	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}
}
