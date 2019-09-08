# Notes

## Items

* Need to break object up into Tika (Geo in the future), break out graphloader and package builder code

## Build sequence for CSDCO Semantic Network

Note:  We may perform a graph drop or other process if we want to cleanly initialize
the graph.  Otherwise we will use specific drop and loads to perform updates along 
the way.

1. (Why do this if step 4?) Load the CSDCOProject.nt and objectGraph.nq from OCD Semantic Network  (Load data dictionary without instances)
2. VaultWalker runs to load objects and object metadata kernel  x-do and x-do-meta respectively
3. GraphBuilder is run to build out "proj" and "borehole" to x-do-resources
4. ObjectEngine>GraphLoader loads x-do-meta and x-do-reosources to graph
5. ObjectEngine>Tika to build the tika index in x-do-tika "Optional"
6. ObjectEngine>PackageBuilder is run to build packages and their metadata in x-do-packages and x-do-packages-meta
7. ObjectEngine>GraphLoader run again to load x-do-tika and x-do-packages-meta

## Dev set up

* Use my Fcore triple store

## Prod set up

## Query

* example DO <http://opencoredata.org/id/do/1d7700ac873af60fe3c7e34c3326671623b6cec649c75aa4e4b303b924624870>

## Reference

```bash
mc cat clear/csdco-do-packages/00fa3009ba0af67da91d9742fd3a64c40877158f7811b2a21c40b90f498ca348
```

```json
{
 "name": "XING",
 "title": "XING",
 "description": "A Frictionless Data Package for XING",
 "sources": [
 {
 "name": "CSDCO",
 "web": "http://csdco.org"
 }
 ],
 "resources": [
 {
 "url": "http://opencoredata.org/id/do/bj6tk07tr9b9lopa8pfg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pg0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pgg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8ph0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8phg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pi0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pig"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pj0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pjg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pk0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pkg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pl0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8plg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ftr9b9lopa8pm0"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk0ntr9b9lopa8pmg"
 },
 {
 "url": "http://opencoredata.org/id/do/bj6tk1ftr9b9lopa8pn0"
 }
 ]
}
```

```bash
mc cat clear/csdco-do-packages-meta/bl4s1oqu6s77r4hp966g 
```

```json
{
 "@context": {
  "@vocab": "http://schema.org/",
  "csdco": "http://opencoredata.org/voc/csdco/1/",
  "re3data": "http://example.org/re3data/0.1/"
 },
 "@id": "http://opencoredata.org/id/do/b5993b8ea246d19e03ae0580d9d5f38921c5b424f519d6e252eb25fac6309e06",
 "@type": "DataSet",
 "about": "http://opencoredata.org/id/do/1001",
 "description": "A project level dataset for CSDCO project 1001",
 "hasPart": [
  {
   "@type": "http://schema.org/DigitalDocument",
   "url": "http://opencoredata.org/id/do/bj6uct7tr9b9lopan07g"
  },
 ...
  {
   "@type": "http://schema.org/DigitalDocument",
   "url": "http://opencoredata.org/id/do/bj6ud3ntr9b9lopan2jg"
  },
  {
   "@type": "http://schema.org/DigitalDocument",
   "url": "http://opencoredata.org/id/do/bj6ud4vtr9b9lopan2l0"
  }
 ],
 "keywords": "scientific drilling geoscience",
 "license": "CC-0",
 "name": "1001 Dataset",
 "url": "http://opencoredata.org/id/do/b5993b8ea246d19e03ae0580d9d5f38921c5b424f519d6e252eb25fac6309e06"
}
```


------------------------------------------------------

Use:

```SPARQL
PREFIX  xsd:    <http://www.w3.org/2001/XMLSchema#>
PREFIX  dc:     <http://purl.org/dc/elements/1.1/>
PREFIX  :       <.>
SELECT DISTINCT ?s ?description ?name ?license ?encodingFormat ?url ?type ?additionType ?dateCreated ?identifier
{
    GRAPH ?g { 
        ?s <http://schema.org/isRelatedTo> "YUFL" .
        ?s <http://schema.org/description> ?description .
        ?s <http://schema.org/name> ?name .
        OPTIONAL { ?s <http://schema.org/license> ?license . }
        OPTIONAL { ?s <http://schema.org/encodingFormat> ?encodingFormat . }
        OPTIONAL { ?s <http://schema.org/url> ?url . }
        OPTIONAL { ?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type . }
        OPTIONAL { ?s <http://schema.org/additionalType> ?additionalType . }
        OPTIONAL { ?s <http://schema.org/dateCreated> ?dateCreated . }
        OPTIONAL { ?s <http://schema.org/identifier> ?identifier . }
    }
}
```

The full FDP resource is
```json
{
"name": "solar-system",
"path": "http://example.com/solar-system.csv",
"title": "The Solar System",
"description": "My favourite data about the solar system.",
"format": "csv",
"mediatype": "text/csv",
"encoding": "utf-8",
"bytes": 1,
"hash": "",
"schema": "",
"sources": "",
"licenses": ""
}
```

These are the parameters we have....
```turtle
<http://schema.org/description>
<http://schema.org/name>
<http://schema.org/license>
<http://schema.org/encodingFormat>
<http://schema.org/url>
<http://www.w3.org/1999/02/22-rdf-syntax-ns#type>
<http://schema.org/additioa;Type>
<http://schema.org/dateCreated>
<http://schema.org/identifier>
<http://schema.org/isRelatedTo>
```

