#ocdWeb

##About
ocdWeb is the web ui for the Open Core Data effort.  It is distinct from ocdServices which is the primary means to query or request data.   ocdWeb implements approaches to explose linked open data, schema.org and W3C CSV patterns.   


##Cross compile and small docker files
Reference: https://sebest.github.io/post/create-a-small-docker-image-for-a-golang-binary/ 
Reference: https://blog.codeship.com/building-minimal-docker-containers-for-go-applications/


###Cross compile the binary (cgo not enabled) 
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go


 /Users/dfils/Data/OCDDataVolumes/webstatic/:/static


### Notes

```
/id/dataset    a dataset
/voc           a vocabulary  (no 303?)
/id/rdf        a resource

303's 

/doc/dataset  	Dataset landing page 
/doc/rdf     	RDF resource landing page
```

```
Fils:ocdWeb dfils$ docker run  -d -p 27017:27017 -v /Users/dfils/Data/OCDDataVolumes/mongo/datadir/:/data/db mongo:latest
Fils:ocdWeb dfils$ export MONGO_HOST=localhost
Fils:ocdWeb dfils$ go run main.go

 docker run  -d -p 27017:27017 -v /Users/dfils/Data/OCDDataVolumes/mongo/datadir/:/data/db mongo:latest

```