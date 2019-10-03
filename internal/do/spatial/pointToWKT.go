package spatial

import (
	"log"
	"strconv"

	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/geojson"
	"github.com/twpayne/go-geom/encoding/wkt"
)

// LatLongWKT convert lat long to WKT
func LatLongWKT(slat, slong string) (string, error) {
	lat, err := strconv.ParseFloat(slat, 64)
	if err != nil {
		return "", err
	}
	long, err := strconv.ParseFloat(slong, 64)
	if err != nil {
		return "", err
	}

	p := geom.NewPoint(geom.XY).MustSetCoords([]float64{long, lat}) // long lat vs lat long
	//log.Println(p)
	s, err := wkt.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	return s, err
}

// LatLongGJS convert lat long to WKT
func LatLongGJS(slat, slong []string) (string, error) {
	g := geom.NewGeometryCollection()

	// TODO..   test len of slat/slong and if 1, don't use a collections
	// just geom.NewGeoemtry with the point...

	for x, _ := range slat {
		lat, err := strconv.ParseFloat(slat[x], 64)
		if err != nil {
			return "", err
		}
		long, err := strconv.ParseFloat(slong[x], 64)
		if err != nil {
			return "", err
		}
		g.MustPush(geom.NewPoint(geom.XY).MustSetCoords(geom.Coord{long, lat}))
	}
	//log.Println(p)
	s, err := geojson.Marshal(g)
	if err != nil {
		log.Println(err)
	}

	return string(s), err
}
