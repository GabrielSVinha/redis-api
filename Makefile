PACKAGE  = hellogopher
BASE     = $(GOPATH)/src/$(PACKAGE)

clean:
	rm -rf redis-api

depend:
	go get -d -v "github.com/gorilla/mux" && go get -d -v "github.com/mediocregopher/radix.v2/redis"

build: depend
	go build redis-api.go
