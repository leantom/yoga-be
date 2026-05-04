.PHONY: run test tidy docker-build deploy-indexes

run:
	go run ./cmd/api

test:
	go test ./...

tidy:
	go mod tidy
	gofmt -w cmd internal

docker-build:
	docker build -t yoga-api:local .

deploy-indexes:
	firebase deploy --only firestore:indexes
