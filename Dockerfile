FROM golang:1.11 as build
ADD . /go/src/github.com/m-lab/metadata-collector
RUN go get -v github.com/m-lab/metadata-collector

FROM alpine:3.7
RUN apk add lshw
COPY --from=build /go/bin/metadata-collector /
WORKDIR /
# Run the collector once to verify that every command it invokes can be invoked inside the container.
RUN /metadata-collector -once
# Remove the created directories to allow them to be mountpoints when deployed.
RUN rm -Rf /var/spool/configuration
RUN rm -Rf /var/spool/hardware
RUN rm -Rf /var/spool/software
# If we made it here, then everything works!