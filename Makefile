VERSION ?= $(shell git describe --tags --always --dirty)

.PHONY: test clean vet fmt fmtcheck build run

bin/ghwc:
	@cd cmd/ghwc && go build -o ../../bin/ghwc main.go && cd ../../

# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
  # use the rest as arguments for "run"
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

build: clean bin/ghwc

run: build
	@bin/ghwc $(RUN_ARGS)

test: vet
	go test -v ./...

fmt:
	@echo "Running gofmt on all sources..."
	@gofmt -s -l -w .

fmtcheck:
	@bash -c "diff -u <(echo -n) <(gofmt -d .)"

vet:
	go vet ./...

clean:
	@rm -f bin/ghwc
