builddist:
	gox -osarch="linux/amd64 linux/arm darwin/amd64 windows/amd64" -output="dist/{{.OS}}/{{.Arch}}/terraform-enterprise-migrator"

test:
	go test -cover
