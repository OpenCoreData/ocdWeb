BINARY := web
DOCKERVER := 0.9.10
.DEFAULT_GOAL := linux

linux:
	cd cmd/web ; \
	GOOS=linux GOARCH=amd64 CG_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build  --tag="opencoredata/ocdweb:latest" --tag="opencoredata/ocdweb:$(DOCKERVER)"  --file=./build/Dockerfile .

