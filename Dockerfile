FROM golang:1.11 as build
ADD . /go/src/github.com/m-lab/nodeinfo
RUN go get -v github.com/m-lab/nodeinfo

FROM alpine:3.7
RUN apk add lshw
COPY --from=build /go/bin/nodeinfo /
WORKDIR /
# Run things once to verify that every command invoked can be invoked inside the container.
RUN /nodeinfo -once
# Remove the created directory to allow it to be a mountpoint when deployed.
RUN rm -Rf /var/spool/nodeinfo
# If we made it here, then everything works!
ENTRYPOINT ["/nodeinfo"]
