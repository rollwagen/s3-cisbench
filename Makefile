Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
# '-s' strip output;  '-w' suppress warnings;
LDFLAGS := "-s -w -X github.com:rollwagen/s3-cisbench/cmd.Version=$(Version) -X github.com:rollwagen/s3-cisbench/cmd.GitCommit=$(GitCommit)"
export GO111MODULE=on
SOURCE_DIRS = cmd internal main.go

.PHONY: all
all: dist

.PHONY: dist
dist:
	mkdir -p bin/
	rm -rf bin/s3-cisbench*
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/s3-cisbench-linux
	GOARCH=arm64 CGO_ENABLED=0 GOOS=linux go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/s3-cisbench-linux-arm64
	CGO_ENABLED=0 GOOS=darwin go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/s3-cisbench-darwin
	GOARCH=arm64 CGO_ENABLED=0 GOOS=darwin go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/s3-cisbench-darwin-m1
	GOOS=windows CGO_ENABLED=0 go build -a -ldflags $(LDFLAGS) -installsuffix cgo -o bin/s3-cisbench.exe
