FROM golang:1.19 AS go
WORKDIR /go/src/goblue
COPY . .
RUN go get -d -v ./...
RUN make

FROM scratch
COPY --from=go /go/src/goblue/bin/goblue /
CMD [ "/goblue" ]
