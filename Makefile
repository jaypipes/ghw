PKGS := $(shell go list ./... | grep -v /vendor)
BIN_DIR := $(GOPATH)/bin
DEP := $(BIN_DIR)/dep
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test
test: dep
	go test $(PKGS)

$(DEP):
	go get -u github.com/golang/dep/cmd/dep

.PHONY: dep
dep: $(DEP)
	$(DEP) ensure
	(cd ghwc/ && $(DEP) ensure)

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	$(GOMETALINTER) --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor

.PHONY: cover
cover:
	$(shell [ -e coverage.out ] && rm coverage.out)
	@echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PKGS),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	@go tool cover -html=coverage-all.out -o=coverage-all.html
