builddist:
	gox -output="dist/{{.OS}}/{{.Arch}}/terraform-enterprise-migrator"

test:
	go test -cover
