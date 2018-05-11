PKGS := $(shell go list ./... | grep -v /vendor)

BIN_DIR := $(GOPATH)/bin
GOVENDOR := $(BIN_DIR)/govendor

$(GOVENDOR):
	go get -u github.com/kardianos/govendor
	govendor

ghwc/vendor/vendor.json: $(GOVENDOR)
	cd ghwc; govendor sync

.PHONY: test
test: ghwc/vendor/vendor.json
	go test $(PKGS)
