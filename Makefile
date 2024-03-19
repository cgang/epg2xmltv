LDFLAGS := "-s -w"

all: build

clean:
	rm -rf bin/

build:
	CGO_ENABLED="0" GOARCH="amd64" go build -ldflags=${LDFLAGS} -o bin/epg2xml cmd/main.go

