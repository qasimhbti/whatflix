all: build test clean

build: build-whatflix

build-whatflix:
	go build -v -i -o buildtest/logservice ./logservice
	go build -v -i -o buildtest/cacheservice ./cacheservice
	go build -v -i -o buildtest/loadbalancer ./loadbalancer
	go build -v -i -o buildtest/web ./web

test:
	go test -v ./pkg/envutils
	go test -v ./pkg/httperrors
	go test -v ./middleware
	go test -v ./internal/jwtpkg
	go test -v ./internal/version
	go test -v ./controller

clean:
	rm -rf buildtest

.PHONY: build build-whatflix test clean