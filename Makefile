export PATH := ./build:$(PATH)
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export CGO_LDFLAGS := -L./
export LD_LIBRARY_PATH := ./

gencov: 
	cd gen && go test -v -coverprofile=c.out
	go tool cover -html=gen/c.out

testlua: build
	cd genlua && ./genlua _client.go
	#cd genlua && ./genlua _rcached.go

test1: 
	cd test && go build -ldflags "-s -w" -tags lua54

build: 
	cd genlua && go build
