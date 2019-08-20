BINARY := ocd_server
DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := linux

linux:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build  --tag="opencoredata/ocdweb:$(DOCKERVER)"  --file=./build/Dockerfile .

dockerlatest:
	docker build  --tag="opencoredata/ocdweb:latest"  --file=./build/Dockerfile .