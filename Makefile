default: install

.PHONY: docs install lint testacc

docs:
	go generate ./...

install:
	go install .

lint:
	golangci-lint run

testacc:
	TF_ACC=1 go test -count=1 -parallel=4 -timeout 10m -v ./...
