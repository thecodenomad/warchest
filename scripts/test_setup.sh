#!/bin/bash

DEFAULT_NETWORK_NAME="warchest-net"
WC_CONTAINER_NAME="warchest"
TEST_CONTAINER_NAME="infantry"

function build_test_container() {
  echo "Building test container - ${TOPDIR}"
  docker build -t ${TEST_CONTAINER_NAME} --file ${TOPDIR}/docker/Dockerfile.python . &> /dev/null
}

# TODO: Add ebash!
function create_docker_network() {
  echo "Creating warchest-net network"
  docker network create --driver bridge "${DEFAULT_NETWORK_NAME}" &> /dev/null
  echo "Network ${DEFAULT_NETWORK_NAME} created"
}

function remove_docker_network() {
  echo "Removing warchest-net network"
  docker network remove "${DEFAULT_NETWORK_NAME}" &> /dev/null
  echo "Network ${DEFAULT_NETWORK_NAME} removed"
}

function start_warchest_container() {
  docker run -d --name "${WC_CONTAINER_NAME}" --network "${DEFAULT_NETWORK_NAME}" --publish 8080:8080 "warchest:latest" &> /dev/null
}

function cleanup(){
  infantry_id=$(docker ps -q --no-trunc --format="{{.ID}}" --filter "name=${TEST_CONTAINER_NAME}")
  warchest_id=$(docker ps -q --no-trunc --format="{{.ID}}" --filter "name=${WC_CONTAINER_NAME}")

  if [ -n "${infantry_id}" ]; then
    docker container stop ${TEST_CONTAINER_NAME} &> /dev/null
    docker container rm ${TEST_CONTAINER_NAME} &> /dev/null
  else
    docker container rm ${TEST_CONTAINER_NAME} &> /dev/null
  fi

  if [ -n "${warchest_id}" ]; then
    docker container stop ${WC_CONTAINER_NAME} &> /dev/null
    docker container rm ${WC_CONTAINER_NAME} &> /dev/null
  else
    docker container rm ${WC_CONTAINER_NAME} &> /dev/null
  fi

  remove_docker_network
}

# Make sure there is a clean state
cleanup

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