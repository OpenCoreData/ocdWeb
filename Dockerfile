# Start from scratch image and add in a precompiled binary
# CGO_ENABLED=0 env GOOS=linux go build .
# docker build  --tag="opencoredata/ocdweb:0.9.4"  .
# docker run -d -p 9900:9900  opencoredata/ocdweb:0.9.4
FROM alpine

# Add in the static elements (could also mount these from local filesystem)
ADD ocdWeb /
# ADD ./static  /static   # Replace with -v mounting the static directory
ADD ./templates  /templates

# Add our binary
CMD ["/ocdWeb"]

# Document that the service listens on this port
EXPOSE 9900
