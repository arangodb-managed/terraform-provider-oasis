SHELL = bash
PROJECT := terraform-provider-oasis

COMMIT := $(shell zutano repo build)
VERSION := $(shell zutano repo version)
DOCKERIMAGE ?= $(shell zutano docker image --name=$(PROJECT))

all: binaries check

clean:
	rm -Rf bin

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm darwin/amd64" \
		-ldflags="-X main.projectVersion=${VERSION} -X main.projectBuild=${COMMIT}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)" \
		-tags="netgo" \
		./...

check:
	zutano go check ./...

.PHONY: test
test:
	mkdir -p bin/test
	go test -coverprofile=bin/test/coverage.out -v ./... | tee bin/test/test-output.txt ; exit "$${PIPESTATUS[0]}"
	cat bin/test/test-output.txt | go-junit-report > bin/test/unit-tests.xml
	go tool cover -html=bin/test/coverage.out -o bin/test/coverage.html

.PHONY: test-acc
test-acc:
	TF_ACC=1 go test -v ./...

bootstrap:
	go get github.com/arangodb-managed/zutano
	go get github.com/mitchellh/gox
	go get github.com/jstemmer/go-junit-report
	go get github.com/hashicorp/terraform-plugin-sdk/plugin
	go get github.com/hashicorp/terraform-plugin-sdk/terraform	

docker:
	docker build \
		--build-arg=GOARCH=amd64 \
		-t $(DOCKERIMAGE) .

docker-push:
	docker push $(DOCKERIMAGE)

.PHONY: update-modules
update-modules:
	zutano update-check --quiet --fail
	test -f go.mod || go mod init
	go mod edit \
		$(shell zutano go mod replacements)
	go get \
		$(shell zutano go mod latest \
			github.com/arangodb-managed/apis \
		)
	go mod tidy
