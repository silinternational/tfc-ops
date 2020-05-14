builddist:
	gox -osarch="linux/amd64 linux/arm darwin/amd64 win/amd64" -output="dist/{{.OS}}/{{.Arch}}/terraform-enterprise-migrator"

test:
	go test -cover
