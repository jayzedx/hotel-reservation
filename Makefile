build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

test-seed:
	@export CONFIG_PATH=./ && go run scripts/test_seed.go
#	 @set CONFIG_PATH=./ && go run scripts/test_seed.go
test:
	@go test -v ./...