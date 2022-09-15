#  Build the binary
#
FROM golang:1.19 AS go
WORKDIR /go/src/goblue
COPY . .
RUN go get -d -v ./...
RUN make
RUN strip bin/goblue

# Create the container,  will contain a sample file in /data, normal use would be to volume mount over /data
#
FROM scratch
COPY --from=go /go/src/goblue/bin/goblue /
COPY pkg/blueheaders/sample.tmp /data/sample.tmp
COPY README.md /
ENTRYPOINT [ "/goblue" ]
CMD ["-p", "9580", "-d", "/data"]
