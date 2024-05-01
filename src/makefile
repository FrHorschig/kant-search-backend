.PHONY: build

test:
	go generate ./...
	go test ./... -tags=integration,unit

int:
	go generate ./...
	go test ./... -tags=integration

unit:
	go generate ./...
	go test ./... -tags=unit

cover:
	go generate ./...
	go test ./... -tags=integration,unit -coverprofile coverage.out && go tool cover -html=coverage.out
