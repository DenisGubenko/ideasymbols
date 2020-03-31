# ideasymbols

## Setup for deployment

````bash
$ go get github.com/DenisGubenko/ideasymbols
$ cd $GOPATH/src/github.com/DenisGubenko/ideasymbols
$ cp configs/.storage.env .
$ cp configs/.http-server.env .
$ make deps
$ make dev_up

Into browser:
http://localhost:1234/request
http://localhost:1234/admin/requests
````

## Setup for development

````bash
$ go get github.com/DenisGubenko/ideasymbols
$ cd $GOPATH/src/github.com/DenisGubenko/ideasymbols
$ make deps
````

