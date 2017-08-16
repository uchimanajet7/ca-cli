BINARY := ca
VERSION	:= v0.1.1
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
deps: ## Get dependent package with go get.
	@go version
	go get -v -t -d ./...
# need windows/amd64
	go get -v github.com/inconshreveable/mousetrap

.PHONY: build
build: deps ## Build a binary that runs in the current environment.(make default)
	go build -v -ldflags $(LDFLAGS) -o ./$(BINARY)

.PHONY: install
install: deps ## Execute "go install" command.
	go install -v -ldflags $(LDFLAGS) -o ./$(BINARY)

.PHONY: cross
cross: deps $(GOX) ## Run cross-compilation and create darwin, linux, windows binaries.
	rm -rf ./out && \
	gox -verbose -ldflags $(LDFLAGS) -osarch $(OSARCH) -output "./out/${GITHUB_REPO}_${VERSION}_{{.OS}}_{{.Arch}}/$(BINARY)"

.PHONY: package
package: cross $(ARCHIVER) ## Run cross-compiling and create release packages for darwin, linux, windows.
	rm -rf ./pkg && mkdir ./pkg && \
	cd out && \
	find * -name "ca*" -type f | awk -F '/' '{system("archiver make ../pkg/" $$1 ".zip " $$0)}' && \
	cd ..

.PHONY: release
release: $(GHR) ## Upload release package to GitHub.
	ghr -delete -u $(GITHUB_ORG) $(VERSION) pkg/

.PHONY: pre-release
pre-release: $(GHR) ## Upload pre-release package to GitHub.
	ghr -delete -prerelease -u $(GITHUB_ORG) $(VERSION) pkg/

.PHONY: digest
digest: ## Display message digest of release package.
	openssl dgst -sha256 pkg/*.zip

.PHONY: lint
lint: $(LINT) ## Execute the "golint" command.
	@golint ./...

.PHONY: vet
vet: ## Execute the "go vet" command.
	@go vet ./...

.PHONY: test
test: ## Execute the "go test" command.
	@go test -v ./...

.PHONY: cover
cover: ## Execute the "go test" command and get coverage.
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

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Display help message.
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
