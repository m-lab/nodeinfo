FROM golang:1.18 as build
ENV CGO_ENABLED 0
COPY . /go/src/github.com/m-lab/nodeinfo
WORKDIR /go/src/github.com/m-lab/nodeinfo
RUN apt update && apt reinstall ca-certificates
RUN go install \
    -v \
    -ldflags "-X github.com/m-lab/go/prometheusx.GitShortCommit=$(git log -1 --format=%h)" \
    ./...

FROM alpine:3.7
# Add all binaries that we may want to run that are not in alpine by default.
RUN apk add --no-cache lshw
COPY --from=build /go/bin/nodeinfo /
COPY --from=build /go/src/github.com/m-lab/nodeinfo/api/nodeinfo1.json /var/spool/datatypes/nodeinfo1.json
WORKDIR /
# Make sure /nodeinfo can run (has no missing external dependencies).
RUN /nodeinfo -h 2> /dev/null
ENTRYPOINT ["/nodeinfo"]
