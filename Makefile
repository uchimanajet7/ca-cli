BINARY := ca
VERSION	:= v0.1.0
REVISION := $(shell git rev-parse --short HEAD)
GITHUB_ORG := uchimanajet7
GITHUB_REPO := ca-cli
LDFLAGS := "-X github.com/${GITHUB_ORG}/${GITHUB_REPO}/cmd.version=${VERSION} -X github.com/${GITHUB_ORG}/${GITHUB_REPO}/cmd.revision=${REVISION} -extldflags \"-static\""
OSARCH := "darwin/amd64 linux/amd64 windows/amd64"

ifndef GOBIN
GOBIN := $(shell echo "$${GOPATH%%:*}/bin")
endif

LINT := $(GOBIN)/golint
GOX := $(GOBIN)/gox
ARCHIVER := $(GOBIN)/archiver
GHR := $(GOBIN)/ghr

$(LINT): ; @go get github.com/golang/lint/golint
$(GOX): ; @go get github.com/mitchellh/gox
$(ARCHIVER): ; @go get github.com/mholt/archiver/cmd/archiver
$(GHR): ; @go get github.com/tcnksm/ghr

.DEFAULT_GOAL := build

.PHONY: deps
deps:
	go get -v -t -d ./... && \
	go get -v github.com/inconshreveable/mousetrap 
	# need windows/amd64

.PHONY: build
build: deps
	go build -v -ldflags $(LDFLAGS) -o ./$(BINARY)

.PHONY: install
install: deps
	go install -ldflags $(LDFLAGS)

.PHONY: cross
cross: deps $(GOX)
	rm -rf ./out && \
	gox -verbose -ldflags $(LDFLAGS) -osarch $(OSARCH) -output "./out/${GITHUB_REPO}_${VERSION}_{{.OS}}_{{.Arch}}/$(BINARY)"

.PHONY: package
package: cross $(ARCHIVER)
	rm -rf ./pkg && mkdir ./pkg && \
	cd out && \
	find * -name "ca*" -type f | awk -F '/' '{system("archiver make ../pkg/" $$1 ".zip " $$0)}' && \
	cd ..

.PHONY: release
release: $(GHR)
	ghr -u $(GITHUB_ORG) $(VERSION) pkg/

.PHONY: digest
digest:
	openssl dgst -sha256 pkg/*.zip

.PHONY: lint
lint: $(LINT)
	@golint ./...

.PHONY: vet
vet:
	@go vet ./...

.PHONY: test
test:
	@go test -v ./...

.PHONY: check
check: lint vet test build

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: cover
cover:
	set -e; \
	COVER_MODE=atomic; \
	COVER_FILE=coverage.txt; \
	echo "mode: $$COVER_MODE" > $$COVER_FILE; \
	for d in $$(go list ./... | grep -v vendor); do \
		go test -race -coverprofile=profile.out -covermode=$$COVER_MODE $$d; \
    	if [ -f profile.out ]; then \
        	cat profile.out | tail -n +2 >> $$COVER_FILE; \
			rm profile.out; \
    	fi \
	done; \
	go tool cover -html=$$COVER_FILE