builddist:
	goreleaser release --snapshot --skip-publish

test:
	go test -cover
