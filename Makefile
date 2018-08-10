PKGS := $(shell go list ./... | grep -v /vendor)
BIN_DIR := $(GOPATH)/bin
GOVENDOR := $(BIN_DIR)/govendor
GOMETALINTER := $(BIN_DIR)/gometalinter

.PHONY: test
test: ghwc/vendor/vendor.json
	go test $(PKGS)

$(GOVENDOR):
	go get -u github.com/kardianos/govendor

ghwc/vendor/vendor.json: $(GOVENDOR)
	cd ghwc; govendor sync

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor
