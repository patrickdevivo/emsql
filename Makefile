gotags = "sqlite_vtable,sqlite_json1"

vet:
	go vet -v -tags=$(gotags) ./...

build:
	go build -v -tags=$(gotags)

lint:
	golangci-lint run --build-tags $(gotags)

test:
	go test -v -tags=$(gotags) ./...

test-cover:
	go test -v -tags=$(gotags) ./... -cover -covermode=count -coverprofile=coverage.out
	go tool cover -html=coverage.out

bench:
	go test -v -tags=$(gotags) -bench=. -benchmem -run=^nomatch ./...