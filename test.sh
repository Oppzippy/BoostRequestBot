#!/bin/bash
docker-compose -f ./docker-integration/docker-compose.yml build
docker-compose -f ./docker-integration/docker-compose.yml up --abort-on-container-exit
