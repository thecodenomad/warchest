#
# Copyright 2021, Ray Gomez <codenomad@gmail.com>
#

#############################
# Essential Build Variables #
#############################
# WARCHEST_CONFIG="${WARCHEST_CONFIG:=$(pwd)/src/config/testdata/CoinConfig.json}"

TOPDIR := ${CURDIR}
CB_API_KEY := ${CB_API_KEY}
CB_API_SECRET := ${CB_API_SECRET}

#############################
# Make Targets              #
#############################

all: clean test build

build:
	go build -o warchest src/main.go

clean:
	rm -rf coverage.out coverage.html warchest

covtest: clean
	go test -v -coverprofile=coverage.out ./...

covreport: clean covtest
	go tool cover -html=coverage.out -o coverage.html

docker: build
	# Populate submodule
	git submodule update --init --remote warchest-ui
	# Build the bits
	docker build . -f Dockerfile --tag warchest

run: docker
	docker run -v ${TOPDIR}/logs:/code/logs --env CB_API_KEY=${CB_API_KEY} --env CB_API_SECRET=${CB_API_SECRET} -p 8080:8080 warchest:latest

demo: docker
	docker run -v ${TOPDIR}/logs:/code/logs --env CB_API_KEY=demo -p 8080:8080 warchest:latest

test:
	go test ./... -v

L2: docker
	${TOPDIR}/scripts/test_setup.sh
