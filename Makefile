.PHONY: test
test: vet
	go test -v ./...

.PHONY: fmt
fmt:
	@echo "Running gofmt on all sources..."
	@gofmt -s -l -w .

.PHONY: fmtcheck
fmtcheck:
	@bash -c "diff -u <(echo -n) <(gofmt -d .)"

.PHONY: vet
vet:
	go vet ./...
