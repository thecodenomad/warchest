#
# Copyright 2021, Ray Gomez <codenomad@gmail.com>
#

#############################
# Essential Build Variables #
#############################
# WARCHEST_CONFIG="${WARCHEST_CONFIG:=$(pwd)/src/config/testdata/CoinConfig.json}"

#############################
# Make Targets              #
#############################

all: clean test build deploy

build:
	go build -o warchest src/main.go

clean:
	rm -rf coverage.out coverage.html warchest

covtest: clean
	go test -v -coverprofile=coverage.out ./...

covreport: clean covtest
	go tool cover -html=coverage.out -o coverage.html

deploy: build
	# docker status

test:
	go test ./... -v
