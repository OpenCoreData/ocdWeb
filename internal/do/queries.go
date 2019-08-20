package do

const queries = `
# Comments are ignored, except those tagging a query.

# tag: test
SELECT ?s ?name ?desc
WHERE {
  ?s ?p <http://opencoredata.org/id/do/cc7481953cacce428eda4f3ed11c96a4ea3b1114084acf29496c15908cb6dee4> .
  ?s <http://schema.org/name> ?name .
  ?s <http://schema.org/description> ?desc
}

# tag: typecheck
SELECT DISTINCT ?type ?graph
WHERE
{  GRAPH ?graph
    {
        <http://opencoredata.org/id/do/{{.OID}}> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
    }
}

# tag: projfeatures
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
    BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	 {
	    ?res <http://schema.org/about> ?ss .
	    ?res <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://opencoredata.org/voc/csdco/v1/Borehole> .
	 }
  }
}

# tag: projdatasets
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
     BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	  {
	     ?res <http://schema.org/about> ?ss .
	     ?res <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataSet> .
	  }
  }
}

# tag: projFDPs
SELECT DISTINCT ?res
WHERE { GRAPH ?graph
{
     BIND ("http://opencoredata.org/id/do/{{.ID}}" AS ?ss)
	  {
		?s <http://schema.org/about> ?ss .
		?s <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataSet> .
		?s <http://schema.org/distribution> ?dist .
		?dist <http://schema.org/contentUrl> ?res
	}
  }
}

# tag: getObject
SELECT  DISTINCT ?date ?mimetype ?type ?license ?filetype ?name ?desc ?related ?url ?text
WHERE {
  ?s ?p "{{.OID}}" .
  ?s2 ?p2 ?s .
  ?s2 <http://schema.org/dateCreated> ?date .
  ?s2 <http://schema.org/encodingFormat> ?mimetype .
  ?s2 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> ?type .
  ?s2 <http://schema.org/license> ?license .
  ?s2 <http://schema.org/additionType> ?filetype .
  ?s2 <http://schema.org/name> ?name .
  ?s2 <http://schema.org/text> ?text .
  ?s2 <http://schema.org/description> ?desc .
  ?s2 <http://schema.org/isRelatedTo> ?related .
  ?s2 <http://schema.org/url> ?url
}
`
