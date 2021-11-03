#!/bin/bash

DEFAULT_NETWORK_NAME="warchest-net"
WC_CONTAINER_NAME="warchest"
TEST_CONTAINER_NAME="infantry"

function build_test_container() {
  echo "Building test container - ${TOPDIR}"
  docker build -t ${TEST_CONTAINER_NAME} --file ${TOPDIR}/docker/Dockerfile.python .
}

# TODO: Add ebash!
function create_docker_network() {
  echo "Creating warchest-net network"
  docker network create --driver bridge "${DEFAULT_NETWORK_NAME}" && echo "Network ${DEFAULT_NETWORK_NAME} created"
}

function remove_docker_network() {
  echo "Removing warchest-net network"
  docker network remove "${DEFAULT_NETWORK_NAME}" && echo "Network ${DEFAULT_NETWORK_NAME} removed"
}

function start_warchest_container() {
  docker run -d --name "${WC_CONTAINER_NAME}" --network "${DEFAULT_NETWORK_NAME}" --publish 8080:8080 "warchest:latest"
}

function cleanup(){
  cypress_id=$(docker ps -q --no-trunc --format="{{.ID}}" --filter "name=${TEST_CONTAINER_NAME}")
  warchest_id=$(docker ps -q --no-trunc --format="{{.ID}}" --filter "name=${WC_CONTAINER_NAME}")

  if [ -n "${cypress_id}" ]; then
    echo "Stopping Cypress container"
    docker container stop ${cypress_id}
    docker container rm ${cypress_id}
  else
    echo "Cypress container not found!"
  fi

  if [ -n "${warchest_id}" ]; then
    echo "Stopping Warchest container"
    docker container stop ${warchest_id}
    docker container rm ${warchest_id}
  else
    echo "Warchest container not found!"
  fi

  remove_docker_network
}

# Setup docker network
create_docker_network

# Setup server in network
start_warchest_container

# Setup python test container in network
build_test_container

# Execute tests
docker run --name ${TEST_CONTAINER_NAME} ${TEST_CONTAINER_NAME}:latest

# Cleanup after testing
cleanup