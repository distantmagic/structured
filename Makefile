.DEFAULT_GOAL := structured

GO_SOURCES := \
	$(wildcard *.go) \
	$(wildcard */*.go)

structured: $(GO_SOURCES)
	go build -o structured

.PHONY: clean
clean:
	rm -f log.db
	rm -f structured
	rm -rf snapshots
	rm -f stable.db

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run
