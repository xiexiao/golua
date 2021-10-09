export PATH := ./build:$(PATH)
export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct
export CGO_LDFLAGS :=  -L$(CURDIR)/ -L./
export LD_LIBRARY_PATH := ./

gencov: 
	cd gen && go test -v -coverprofile=c.out
	go tool cover -html=gen/c.out

fmt: 
	cd lua && go fmt 
	cd gen && go fmt 
	cd genlua && go fmt 

testlua1: build
	cd genlua && ./genlua _client.go
	#cd genlua && ./genlua _rcached.go

testlua2:
	cd test && go build -ldflags "-s -w" -tags lua54

test1: 
	cd gen && go test -v -test.run TestGenBytes

testlua:
	go test ./lua -v -ldflags "-s -w" -tags lua54

t1:
	go test ./lua -v -ldflags "-s -w" -tags lua54 -test.run TestLPushToTable

t2:
	#go test ./gen -v -ldflags "-s -w" -tags lua54 -test.run TestSetMetatable
	cd gen && go test -v -test.run TestGenNew
	#cd gen && go test -v 

t3:
	cd gen && go test -v -test.run TestGenUseSame

t4:
	cd gen && go test -v -test.run TestGenLuaFunc

t5:
	cd gen && go test -v -test.run TestGenLuaMT

t6:
	cd gen && go test -v -test.run TestGenSlice

main1: 
	cd main && go build -ldflags "-s -w" -tags lua54

build: 
	cd genlua && go build

install: 
	cd genlua && go install
