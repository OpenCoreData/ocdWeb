# Start from scratch image and add in a precompiled binary
# CGO_ENABLED=0 env  go build .
# docker build  --tag="opencoredata/ocdweb:0.9.4"  .
# docker run -d -p 9900:9900  opencoredata/ocdweb:0.9.4
#FROM alpine
FROM scratch

# Add in the static elements (could also mount these from local filesystem)
ADD cmd/web/web  /

# ADD ../web/static  /static   # Replace with -v mounting the /web/static directory
ADD web/templates  /templates

# Add our binary
CMD ["/web"]

# Document that the service listens on this port
EXPOSE 9900
