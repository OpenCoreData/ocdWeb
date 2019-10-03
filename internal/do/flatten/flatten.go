package flatten

import (
	"encoding/json"
	"log"

	"github.com/piprate/json-gold/ld"
)

// TestFlatten to flatten the doc
func TestFlatten(jsonld string) string {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	var myInterface interface{}
	err := json.Unmarshal([]byte(jsonld), &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
	}

	flattenedDoc, err := proc.Flatten(myInterface, nil, options)

	ld.PrintDocument("JSON-LD graph section", flattenedDoc) // debug print....

	jsonm, err := json.MarshalIndent(flattenedDoc, "", " ")
	if err != nil {
		log.Println("Error trying to marshal data", err)
	}

	return string(jsonm)

}
