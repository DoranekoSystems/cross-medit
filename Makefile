GOCMD=go
GOTEST=$(GOCMD) test -v
GOBUILD=$(GOCMD) build
BINARY_NAME=medit

all: build deploy

test:
	$(GOTEST) ./pkg/*

build-windows-amd64:
	GOOS=windows \
	GOARCH=amd64 \
	$(GOBUILD) -o $(BINARY_NAME).exe

build-ios-arm64:
	CGO_ENABLED=1 \
	GOOS=ios \
	GOARCH=arm64 \
	SDK=iphoneos \
	CC=$(PWD)/clangwrap.sh \
	$(GOBUILD) -o $(BINARY_NAME)

clean:
	rm $(BINARY_NAME)
