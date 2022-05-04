DIR?=./autorest/

default: build

build: fmt
	cd $(DIR); go install

test:
	cd $(DIR); go test -v

vet:
	cd $(DIR); go vet

fmt:
	gofmt -w $(DIR)

.PHONY: build test vet fmt
