FROM golang:1.12 as build
ENV CGO_ENABLED 0
COPY . /go/src/github.com/m-lab/nodeinfo
WORKDIR /go/src/github.com/m-lab/nodeinfo
RUN go get \
    -v \
    -ldflags "-X github.com/m-lab/go/prometheusx.GitShortCommit=$(git log -1 --format=%h)" \
    ./...

FROM alpine:3.7
# Add all binaries that we may want to run that are not in alpine by default.
RUN apk add --no-cache lshw
COPY --from=build /go/bin/nodeinfo /
WORKDIR /
ENTRYPOINT ["/nodeinfo"]
