#!/bin/bash
make full-build
docker build -t angular2guy/reactandgo:latest --build-arg APP_FILE=react-and-go --no-cache .
docker run -p 8080:8080 --memory="192m" --network="host" angular2guy/reactandgo:latest