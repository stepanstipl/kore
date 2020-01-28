ARG GOVERSION
FROM golang:${GOVERSION} as with-source
COPY . .
RUN make go-swagger

# Clean build container with just swagger and jq
FROM golang:${GOVERSION}

RUN apt update -y && apt install jq -y
COPY --from=with-source /go/bin/swagger /go/bin/swagger
CMD bash
