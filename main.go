package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"opencoredata.org/ocdWeb/catalogs"
	"opencoredata.org/ocdWeb/colls"
	"opencoredata.org/ocdWeb/doc"
	"opencoredata.org/ocdWeb/dx"
	"opencoredata.org/ocdWeb/rx"
	"opencoredata.org/ocdWeb/services"
	"opencoredata.org/ocdWeb/voc"
	// _ "net/http/pprof"
)

type MyServer struct {
	r *mux.Router
}

func main() {
	// Common files like; css, js, images, etc...
	rcommon := mux.NewRouter()
	rcommon.PathPrefix("/common/").Handler(http.StripPrefix("/common/", http.FileServer(http.Dir("./static"))))
	http.Handle("/common/", &MyServer{rcommon})

	// root
	parking := mux.NewRouter()
	parking.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", &MyServer{parking})

	// Simpler services to support the web UI  (other services in ocdService)
	servroute := mux.NewRouter()
	servroute.HandleFunc("/services/grid", services.GetGrid)
	servroute.HandleFunc("/services/csdco", services.CSDCOGrid)
	http.Handle("/services/", servroute)

	// Recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// Example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/graph/{id}", dx.RDFRedirection)
	dxroute.HandleFunc("/id/graph/{id}/provenance", dx.ProvRedirection)   // PROV: prov redirection
	dxroute.HandleFunc("/id/graph/{id}/pingback", dx.PingbackRedirection) // PROV: pingback for this resource  (would prefer a master /prov or server)
	dxroute.HandleFunc("/id/dataset/{UUID}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}/{SITE}/{HOLE}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}/{SITE}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}", dx.Redirection)
	dxroute.HandleFunc(`/id/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, dx.Redirection)
	http.Handle("/id/", dxroute)

	// MD5 concept from indie web thoughts...
	// psuedo code == dxroute.HandleFunc("/id/md5/{md5hash}, dx.MD5Redirection")

	// Some early Prov Pingback work here...   Deal with void...  (show void..  allow .rdf file downloads)
	rdfdocs := mux.NewRouter()
	rdfdocs.HandleFunc("/rdf/graph/{id}", rx.RenderWithProvHeader)      // PROV: test cast with Void..  would need to generalize
	rdfdocs.HandleFunc("/rdf/graph/{id}/provenance", rx.RenderWithProv) // PROV: test cast with Void..  would need to generalize
	rdfdocs.HandleFunc("/rdf/graph/{id}/pingback", rx.ProvPingback)     // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/rdf/", rdfdocs)

	// Display Vocabulary entries.  A simple human view.. For machines, check for accepts headers?  no 303 for these?
	vocroute := mux.NewRouter()
	vocroute.PathPrefix("/voc/1/ocdSKOS.ttl").Handler(http.StripPrefix("/voc/", http.FileServer(http.Dir("./static/voc"))))
	vocroute.PathPrefix("/voc/janus/1/ocdJanusSKOS.ttl").Handler(http.StripPrefix("/voc/janus", http.FileServer(http.Dir("./static/voc/janus"))))
	vocroute.HandleFunc("/voc/janus/{version}/{term}", voc.VocJanus)
	vocroute.HandleFunc("/voc/{version}/{term}", voc.VocCore)
	http.Handle("/voc/", vocroute)

	//Browser by id redirection to doc  (gets a specific dataset)  http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B
	docroute := mux.NewRouter()
	docroute.HandleFunc("/doc/dataset/{UUID}", doc.UUIDRender)
	docroute.HandleFunc("/doc/expedition/{LEG}/{SITE}/{HOLE}", doc.ShowFeature)
	docroute.HandleFunc("/doc/expedition/{LEG}/{SITE}", doc.ShowFeature)
	docroute.HandleFunc("/doc/expedition/{LEG}", doc.ShowExpedition)
	docroute.HandleFunc(`/doc/resource/people/{resourcepath:[a-zA-Z0-9=\-\/]+}`, doc.PersonResourceRender) // for GeoLink All Hands Demo, remove afterwards, dont' want person specific version
	docroute.HandleFunc(`/doc/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, doc.ResourceRender)
	docroute.HandleFunc("/doc/dataset/{measurement}/{leg}/{site}/{hole}", doc.Render)
	http.Handle("/doc/", docroute)

	//  Should this catalog?
	// Browse by collection   measurement leg site hole
	// Later Browse options might include:  units, observations. geologic time
	// TODO  worry about namespace collision here...  (need operator ID ?)
	collections := mux.NewRouter()
	// collections.HandleFunc("/collections", colls.Landing)
	// Looking to add in a master catalog collection....   perhaps here or down in the catalog router..  (need to clean both these up)
	collections.HandleFunc("/collections/catalogs", colls.Catalogs)
	collections.HandleFunc("/collections/matrix", colls.MLCounts)          //  IODP matrix
	collections.HandleFunc("/collections/expeditions", doc.AllExpeditions) // Big list view
	// collections.HandleFunc("/collections/expeditions/{LEG}", doc.ShowExpedition)
	// collections.HandleFunc("/collections/januslegs", colls.JanusLegs)
	// collections.HandleFunc("/collections/janusmeasurements", colls.JanusMeasurements)
	collections.HandleFunc("/collections/csdco", colls.CSDCOOverview)                        // CSDCO Matrix
	collections.HandleFunc("/collections/csdco/{HoleID}", colls.CSDCOcollection)             //  landing page for collection of files with a HoleID
	collections.HandleFunc("/collections/csdco/project/{ProjectID}", colls.CSDCOProjectInfo) //  landing page for CSDCO Project information
	collections.HandleFunc("/collections/measurement/{measurements}/{leg}", colls.MLURLSets) //  called from the jrso matrix page
	collections.HandleFunc("/collections/measurement/{measurements}", colls.MesSets)
	// collections.HandleFunc("/collections/leg/{leg}", colls.LegSets)  DEPRECTATED for /doc/expedition/{leg}
	http.Handle("/collections/", collections)

	// Catalog handler..   perhaps a bit redundant with the collections handler above..  need to review this
	catalog := mux.NewRouter()
	catalog.HandleFunc("/catalog/geolink", catalogs.GeoLinkCatalog)
	catalog.HandleFunc("/catalog/geolink/dataset/{resourcepath}", catalogs.GeoLinkDataset)
	// catalog.HandleFunc(`/dataset/geolink/{resourcepath:[a-zA-Z0-9=\-\/]+}`, catalogs.GeoLinkDataset)
	http.Handle("/catalog/", catalog)

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")
	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// ref http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop here if its Preflighted OPTIONS request
	// if req.Method == "OPTIONS" {
	// 	return
	// }

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
