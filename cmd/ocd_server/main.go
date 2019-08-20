package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
	"opencoredata.org/ocdWeb/internal/colls"
	"opencoredata.org/ocdWeb/internal/datapkg"
	"opencoredata.org/ocdWeb/internal/do"
	"opencoredata.org/ocdWeb/internal/doc"
	"opencoredata.org/ocdWeb/internal/dx"
	"opencoredata.org/ocdWeb/internal/services"
	"opencoredata.org/ocdWeb/internal/utils"
	"opencoredata.org/ocdWeb/internal/voc"
	// _ "net/http/pprof"
)

// MyServer is the Gorilla mux router structure
type MyServer struct {
	r *mux.Router
}

var minioVal, portVal, accessVal, secretVal, bucketVal string
var sslVal bool

func init() {
	akey := os.Getenv("MINIO_ACCESS_KEY")
	skey := os.Getenv("MINIO_SECRET_KEY")
	mhost := os.Getenv("MINIO_HOST")
	mport := os.Getenv("MINIO_PORT")

	flag.StringVar(&minioVal, "address", mhost, "FQDN for server")
	flag.StringVar(&portVal, "port", mport, "Port for minio server, default 9000")
	flag.StringVar(&accessVal, "access", akey, "Access Key ID")
	flag.StringVar(&secretVal, "secret", skey, "Secret access key")
	flag.StringVar(&bucketVal, "bucket", "provisium", "The configuration bucket")
	flag.BoolVar(&sslVal, "ssl", false, "Use SSL boolean")
}

func minioHandler(minioClient *minio.Client,
	f func(minioClient *minio.Client, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(minioClient, w, r) })
}

func main() {
	// Load configurations
	flag.Parse()
	//minioClient := utils.MinioConnection(minioVal, portVal, accessVal, secretVal, bucketVal, sslVal)
	mc := utils.MinioConnectionDEV() // minio connection

	// Common assets like; css, js, images, etc...
	rcommon := mux.NewRouter()
	rcommon.PathPrefix("/common/").Handler(http.StripPrefix("/common/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/common/", &MyServer{rcommon})

	// root
	parking := mux.NewRouter()
	parking.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/static"))))
	http.Handle("/", &MyServer{parking})

	csdco := mux.NewRouter()
	csdco.PathPrefix("/org/csdco").Handler(http.StripPrefix("/org/csdco/", http.FileServer(http.Dir("./web/static/csdco"))))

	// TODO..   review all this handler?   pointless now?
	//csdco.HandleFunc("/org/csdco/{ID}", do.ObjectView)
	http.Handle("/org/csdco/", csdco)

	// TODO Make a CSDCO router.  Is there really just /csdco name space for pages and resources
	// come from the id and do/doc prefixes?

	// Recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// Example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/graph/{id}", dx.RDFRedirection)
	// dxroute.HandleFunc("/id/graph/{id}/provenance", dx.ProvRedirection)   // PROV: prov redirection
	// dxroute.HandleFunc("/id/graph/{id}/pingback", dx.PingbackRedirection) // PROV: pingback for this resource  (would prefer a master /prov or server)
	dxroute.HandleFunc("/id/dataset/{UUID}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}/{SITE}/{HOLE}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}/{SITE}", dx.Redirection)
	dxroute.HandleFunc("/id/expedition/{LEG}", dx.Redirection)
	dxroute.HandleFunc(`/id/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, dx.Redirection)
	dxroute.HandleFunc(`/id/resource/csdco/feature/{HoleID}`, colls.CSDCOcollection) // DEPRECATE

	dxroute.Handle("/id/do/{ID}", minioHandler(mc, do.ObjectView)).Methods("GET")
	//dxroute.HandleFunc("/id/do/{ID}", do.ObjectView)

	http.Handle("/id/", dxroute)

	//Browser by id redirection to doc  (gets a specific dataset)  http://opencoredata.org/doc/dataset/JanusAgeDatapoint/108/668/B
	docroute := mux.NewRouter()
	docroute.HandleFunc("/doc/dataset/{UUID}", doc.UUIDRender)
	docroute.HandleFunc("/doc/expedition/{LEG}/{SITE}/{HOLE}", doc.ShowFeature)
	docroute.HandleFunc("/doc/expedition/{LEG}/{SITE}", doc.ShowFeature)
	docroute.HandleFunc("/doc/expedition/{LEG}", doc.ShowExpedition)
	docroute.HandleFunc(`/doc/resource/csdco/feature/{HoleID}`, colls.CSDCOcollection)
	docroute.HandleFunc(`/doc/resource/people/{resourcepath:[a-zA-Z0-9=\-\/]+}`, doc.PersonResourceRender) // for GeoLink All Hands Demo, remove afterwards, dont' want person specific version
	docroute.HandleFunc(`/doc/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, doc.ResourceRender)
	docroute.HandleFunc("/doc/dataset/{measurement}/{leg}/{site}/{hole}", doc.Render)

	// TODO  review the fact I never get here !!!!!!!!1
	docroute.Handle("/doc/do/{ID}", minioHandler(mc, do.ObjectView)).Methods("GET")
	//docroute.HandleFunc("/doc/do/{ID}", do.ObjectView)

	docroute.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/doc/", docroute)

	// Collection handler...   elements will be pruned for the DO Cloud approach
	collections := mux.NewRouter()
	collections.HandleFunc("/collections/catalogs", colls.Catalogs)
	collections.HandleFunc("/collections/matrix", colls.MLCounts)          //  IODP matrix
	collections.HandleFunc("/collections/expeditions", doc.AllExpeditions) // Big list view

	// DEPRECATED CSDCO routes to remove
	collections.HandleFunc("/collections/csdco", colls.CSDCOOverview)               // CSDCO Matrix
	collections.HandleFunc("/collections/csdco/{HoleID}", colls.CSDCOcollection)    //  landing page for collection of files with a HoleID
	collections.HandleFunc("/collections/csdco/abstract/{ID}", colls.CSDCOAbstract) // CSDCO abstract display
	// end CSDCO routes to remove

	collections.HandleFunc("/collections/csdco/project/{ProjectID}", colls.CSDCOProjectInfo) //  landing page for CSDCO Project information
	collections.HandleFunc("/collections/measurement/{measurements}/{leg}", colls.MLURLSets) //  called from the jrso matrix page
	collections.HandleFunc("/collections/measurement/{measurements}", colls.MesSets)
	collections.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/collections/", collections)

	// DEPRECATED by the DO Cloud
	// Server Frictionless Data Packages to a Landing Page formed by the schema.org file in the metadata directory
	packages := mux.NewRouter()
	packages.HandleFunc("/pkg/id/{id}.zip", datapkg.DownloadPkg)
	packages.HandleFunc("/pkg/id/{id}", datapkg.ServePkg)
	packages.HandleFunc(`/pkg/id/{id}/{resourcepath:[a-zA-Z0-9=\_\.\-\/]+}`, datapkg.DownloadPkgFile)
	http.Handle("/pkg/", packages)

	// DEPRECATED (will be replaced by LDN and GoLDeN)
	// Some early Prov Pingback work here...   Deal with void...  (show void..  allow .rdf file downloads)
	// rdfdocs := mux.NewRouter()
	// rdfdocs.HandleFunc("/rdf/graph/{id}", rx.RenderWithProvHeader)      // PROV: test cast with Void..  would need to generalize
	// rdfdocs.HandleFunc("/rdf/graph/{id}/provenance", rx.RenderWithProv) // PROV: test cast with Void..  would need to generalize
	// rdfdocs.HandleFunc("/rdf/graph/{id}/pingback", rx.ProvPingback)     // PROV: pingback for this resource  (would prefer a master /prov or server)
	// rdfdocs.NotFoundHandler = http.HandlerFunc(notFound)
	// http.Handle("/rdf/", rdfdocs)

	// TODO..  Should all services be in services?  MOVE ALL THESE TO SERVICES
	// Simpler services to support the web UI  (other services in ocdService)
	servroute := mux.NewRouter()
	servroute.HandleFunc("/services/grid", services.GetGrid)
	servroute.HandleFunc("/services/csdco", services.CSDCOGrid)
	servroute.HandleFunc("/services/csdcov2", services.CSDCOGridv2) // TEST..  remove...
	http.Handle("/services/", servroute)

	// Display Vocabulary entries.  A simple human view.. For machines, check for accepts headers?  no 303 for these?
	vocroute := mux.NewRouter()
	vocroute.PathPrefix("/voc").Handler(http.StripPrefix("/voc/", http.FileServer(http.Dir("./web/static/voc"))))
	vocroute.PathPrefix("/voc/1/ocdSKOS.ttl").Handler(http.StripPrefix("/voc/", http.FileServer(http.Dir("./web/static/voc"))))
	vocroute.PathPrefix("/voc/janus/1/ocdJanusSKOS.ttl").Handler(http.StripPrefix("/voc/janus", http.FileServer(http.Dir("./web/static/voc/janus"))))
	vocroute.HandleFunc("/voc/janus/{version}/{term}", voc.VocJanus)
	vocroute.HandleFunc("/voc/{version}/{term}", voc.VocCore)
	vocroute.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/voc/", vocroute)

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")
	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/404.html", 303)
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
