GOPATH:=$(CURDIR)
export GOPATH

all: output

fmt:
	#gofmt -l -w -s src/

dep:fmt

build:dep
	go build -o bin/gopher main

clean:
	rm -rf output
	rm -rf bin/impr
output:build
	mkdir -p output/bin
	mkdir -p output/conf
	mkdir -p output/log
	mkdir -p output/test
	mkdir -p output/web
	cp -r bin/gopher output/bin/
	cp -r conf/* output/conf/
	cp -r test/* output/test/
	cp -r test/* output/web/

