SHELL = bash
PROJECT := terraform-provider-oasis

COMMIT := $(shell zutano repo build)
VERSION := $(subst v,,$(shell zutano repo version))
DOCKERIMAGE ?= $(shell zutano docker image --name=$(PROJECT))

all: binaries check

clean:
	rm -Rf bin

binaries:
	CGO_ENABLED=0 gox \
		-osarch="linux/amd64 linux/arm linux/arm64 darwin/amd64 darwin/arm64" \
		-ldflags="-X main.projectVersion=${VERSION} -X main.projectBuild=${COMMIT}" \
		-output="bin/{{.OS}}/{{.Arch}}/$(PROJECT)-$(VERSION)" \
		-tags="netgo" \
		./...
	mkdir -p assets
	cd bin/linux/amd64 ; zip -D ../../../assets/$(PROJECT)_$(VERSION)_linux_amd64.zip $(PROJECT)-$(VERSION)
	cd bin/linux/arm64 ; zip -D ../../../assets/$(PROJECT)_$(VERSION)_linux_arm64.zip  $(PROJECT)-$(VERSION)
	cd bin/darwin/amd64 ; zip -D ../../../assets/$(PROJECT)_$(VERSION)_darwin_amd64.zip $(PROJECT)-$(VERSION)
	cd bin/darwin/arm64 ; zip -D ../../../assets/$(PROJECT)_$(VERSION)_darwin_arm64.zip $(PROJECT)-$(VERSION)
	cd assets ; shasum -a 256 *.zip > $(PROJECT)_$(VERSION)_SHA256SUMS
	cp -f terraform-registry-manifest.json assets
	cd assets ; shasum -a 256 *.zip terraform-registry-manifest.json > $(PROJECT)_$(VERSION)_SHA256SUMS

check:
	zutano go check ./...

prepare-release: binaries
	PROJECT=$(PROJECT) VERSION=$(VERSION) ./sign.sh

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
	GOSUMDB=off GOPROXY=direct GO111MODULE=on go get github.com/arangodb-managed/zutano
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
	go get \
		$(shell zutano go mod latest \
			github.com/arangodb-managed/apis \
		)
	go mod tidy
