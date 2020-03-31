FROM golang:1.14-alpine

WORKDIR $GOPATH/src/github.com/DenisGubenko/ideasymbols

COPY . .

RUN apk update && apk add make git gcc libc-dev
RUN make http_service

CMD ["./http_service"]
