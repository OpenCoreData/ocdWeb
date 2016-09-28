# Start from scratch image and add in a precompiled binary
# docker build  --tag="opencoredata/ocdweb:0.1"  .
# docker run -d -p 9900:9900  opencoredata/ocdweb:0.1
FROM alpine

# Add in the static elements (could also mount these from local filesystem)
ADD ocdWeb /
ADD ./static  /static
ADD ./templates  /templates

# Add our binary
CMD ["/ocdWeb"]

# Document that the service listens on this port
EXPOSE 9900
