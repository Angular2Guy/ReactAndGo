#!/bin/bash
export GOGC=off
export GOMEMLIMIT=64MiB
export GODEBUG=gctrace=1
#echo $GOGC $GOMEMLIMIT $GODEBUG
#go build -toolexec=/bin/time
#./react-and-go
go run main.go
