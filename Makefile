export GO111MODULE := on
export GOPROXY := https://goproxy.cn,direct

gencov: 
	cd gen && go test -v -coverprofile=c.out
	go tool cover -html=gen/c.out

testlua: build
	cd genlua && ./genlua _client.go

build: 
	cd genlua && go build
