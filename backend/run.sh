#!/bin/bash
GOGC=off
GOMEMLIMIT=32MiB
GODEBUG=gctrace=1
#echo $GOGC $GOMEMLIMIT $GODEBUG
#go build
#./react-and-go
go run main.go
