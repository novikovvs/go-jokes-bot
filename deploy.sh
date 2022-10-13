#!/bin/bash

docker stop $(docker ps -a -q)

docker-compose up -d --build go
