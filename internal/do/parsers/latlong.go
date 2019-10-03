package parsers

import (
	"log"

	"github.com/tidwall/gjson"
)

// ProjLatLong get the lat and long from a framed result from Project
func ProjLatLong(j string) ([]string, []string, error) {
	// Could looking for an array be more elegant?   Cue ?
	t := gjson.Get(j, "0.location.geo")
	log.Println(len(t.Array()))

	var lat, long []string
	var latres, longres gjson.Result

	if len(t.Array()) < 1 {
		latres = gjson.Get(j, "0.location.geo.latitude") // need to detect is an array and alter the option here
		longres = gjson.Get(j, "0.location.geo.longitude")
	} else {
		latres = gjson.Get(j, "0.location.geo.#.latitude") // need to detect is an array and alter the option here
		longres = gjson.Get(j, "0.location.geo.#.longitude")
	}

	for _, v := range latres.Array() {
		l := v.String()
		lat = append(lat, l)
	}

	for _, v := range longres.Array() {
		l := v.String()
		long = append(long, l)
	}

	return lat, long, nil
}

// DataLatLong get the lat and long from a framed result on type DataSet
func DataLatLong(j string) ([]string, []string, error) {
	// Could looking for an array be more elegant?   Cue ?
	t := gjson.Get(j, "0.spatialCoverage.geo")
	log.Println(len(t.Array()))

	var lat, long []string
	var latres, longres gjson.Result

	if len(t.Array()) < 1 {
		latres = gjson.Get(j, "0.spatialCoverage.geo.latitude") // need to detect is an array and alter the option here
		longres = gjson.Get(j, "0.spatialCoverage.geo.longitude")
	} else {
		latres = gjson.Get(j, "0.spatialCoverage.geo.#.latitude") // need to detect is an array and alter the option here
		longres = gjson.Get(j, "0.spatialCoverage.geo.#.longitude")
	}

	for _, v := range latres.Array() {
		l := v.String()
		lat = append(lat, l)
	}

	for _, v := range longres.Array() {
		l := v.String()
		long = append(long, l)
	}

	return lat, long, nil
}
