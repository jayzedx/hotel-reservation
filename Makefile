build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test-seed:
#	 @export SEED_MODE=TEST && go run scripts/seed.go
	 @set SEED_MODE=TEST && go run scripts/seed.go
test:
	@go test -v ./...