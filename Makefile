.PHONY: all covtest fmt

all: test build deploy

build:
	go build -o warchest src/main.go

clean:
	rm -rf coverage.out coverage.html warchest

covtest:
	go test -coverprofile=coverage.out ./...

covreport: covtest
	go tool cover -html=coverage.html

deploy: build
	docker status

test:
	go test

