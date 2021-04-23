builddist:
	gox -osarch="linux/amd64 linux/arm darwin/amd64 windows/amd64" -output="dist/{{.OS}}/{{.Arch}}/tfc-ops"

test:
	go test -cover
