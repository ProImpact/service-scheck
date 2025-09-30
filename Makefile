
run: build_app
	build/service_check

.PHONY: build
build: build_app
	go build -o build/test_server cmd/test_server/*

build_app:
	go build -o build/service_check cmd/service_check/*

test:
	go test -v -count=1 ./...

check: build
	go run main.go create --path example/individual-service.json 