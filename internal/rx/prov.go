package rx

import (
	"bufio"
	"bytes"
	"fmt"
	"log"

	"github.com/knakk/rdf"
)

// getProvRecord an un-exported function to generate MOCK prov data for testing
// This is a test function only
func getProvRecord() string {

	tr := []rdf.Triple{}

	// Add in
	newsub, _ := rdf.NewIRI("http://foo.org/thisSample")
	newpred1, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj1, _ := rdf.NewIRI("http://www.w3.org/ns/prov#entity")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred1, Obj: newobj1})

	ga, _ := rdf.NewIRI("http://opencoredata.org/org") // ?
	newpred2, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj2, _ := rdf.NewIRI("http://www.w3.org/ns/prov#org")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred2, Obj: newobj2})

	bn, _ := rdf.NewBlank("bn1")

	newpred3, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj3, _ := rdf.NewIRI("http://www.w3.org/ns/prov#attribution")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred3, Obj: newobj3})

	newpred4, _ := rdf.NewIRI("http://www.w3.org/ns/prov#agent")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred4, Obj: ga})

	newpred5, _ := rdf.NewIRI("http://www.w3.org/ns/prov#hadRole")
	newobj5, _ := rdf.NewIRI("http://www.aurole.org/Publisher")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred5, Obj: newobj5})

	newpred6, _ := rdf.NewIRI("http://www.w3.org/ns/prov#qualifiedAttribution")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred6, Obj: bn})

	newpred7, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj7, _ := rdf.NewIRI("http://www.w3.org/ns/prov#Agent")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred7, Obj: newobj7})

	newpred8, _ := rdf.NewIRI("http://www.w3.org/ns/prov#wasAttributedTo")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred8, Obj: bn})

	newpred9, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#label")
	newobj9, _ := rdf.NewLiteral("Geoscience Australia")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred9, Obj: newobj9})

	fmt.Println(tr)

	var inoutFormat rdf.Format
	inoutFormat = rdf.Turtle //NTriples

	// Create a buffer io writer
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	enc := rdf.NewTripleEncoder(foo, inoutFormat)
	err := enc.EncodeAll(tr)
	err = enc.Close()
	if err != nil {
		log.Println(err)
	}

	foo.Flush()
	return string(b.Bytes())

}
