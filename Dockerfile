FROM golang:1.12 as build
ENV CGO_ENABLED 0
ADD . /go/src/github.com/m-lab/nodeinfo
WORKDIR /go/src/github.com/m-lab/nodeinfo
RUN go get \
    -v \
    -ldflags "-X github.com/m-lab/go/prometheusx.GitShortCommit=$(git log -1 --format=%h)" \
    ./...

FROM alpine:3.7
RUN apk add lshw
COPY --from=build /go/bin/nodeinfo /
WORKDIR /
# Run things once to verify that every command invoked can be invoked inside the container.
RUN mkdir smoketest && /nodeinfo -smoketest -datadir smoketest && rm -Rf /smoketest
# Remove the created directory to allow it to be a mountpoint when deployed.
RUN rm -Rf /var/spool/nodeinfo
# If we made it here, then everything works!
ENTRYPOINT ["/nodeinfo"]
