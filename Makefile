install-deps:
	go get github.com/urfave/cli
	go get github.com/mattn/go-sqlite3
	go get github.com/gin-gonic/gin
	go get github.com/sirupsen/logrus

install-brewmctl:
	go install -v ./cmd/brewmctl

install-brewmserver:
	go install -v ./cmd/brewmserver

install-all: install-brewmctl install-brewmserver


build-local:
	go build -v -o ./build/brewmctl ./cmd/brewmctl/main.go
	go build -v -o ./build/brewmserver ./cmd/brewmserver/main.go

builder:
	docker build -t brewm-builder -f ./build.Dockerfile .

build-pi:
	docker run --rm -v ${GOPATH}:/go -w ${GOPATH}/src/github.com/brewm/gobrewmmer \
	  -e "GOOS=linux" -e "GOARCH=arm" -e "GOARM=6" -e "CGO_ENABLED=1" brewm-builder:latest \
	  go build -v -o /go/src/github.com/brewm/gobrewmmer/build/pi/brewmserver \
	 /go/src/github.com/brewm/gobrewmmer/cmd/brewmserver/main.go

build-all: clean build-local build-pi


clean:
	go clean
	rm -rf ./build
