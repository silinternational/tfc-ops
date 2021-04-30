help:
	@echo "To build, run \"goreleaser release --snapshot --skip-publish\" or for a full release push a new tag to \"origin\"."

test:
	go test -cover
