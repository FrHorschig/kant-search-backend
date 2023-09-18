.PHONY: build

test:
	go test ./... -tags=integration,unit

int:
	go test ./... -tags=integration

unit:
	go test ./... -tags=unit

cover:
	go test ./... -tags=integration,unit -coverprofile coverage.out && go tool cover -html=coverage.out
