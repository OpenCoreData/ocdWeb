package do

import (
	"log"

	"opencoredata.org/ocdWeb/internal/do/framing"
	"opencoredata.org/ocdWeb/internal/do/parsers"
	"opencoredata.org/ocdWeb/internal/do/spatial"
)

func WKT() string {
	return "The wkt view of the object"
}

// GeoJSON convert the data graph into a GeoJSON representation.
// This is not the do bytestream, but rather the metadata object
func GeoJSON() string {
	// see the work in P418/garden/jldSpatial

	sfr := framing.ProjSpatial("json ld")
	lat, long, err := parsers.ProjLatLong(sfr)
	if err != nil {
		log.Println(err)
	}

	gjs, err := spatial.LatLongGJS(lat, long)
	if err != nil {
		log.Println(err)
	}

	return gjs
}
