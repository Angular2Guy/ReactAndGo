#!/bin/bash
GOGC=off
GOMEMLIMIT=32MiB
GODEBUG=gctrace=1
#echo $GOGC $GOMEMLIMIT $GODEBUG
#go build -toolexec=/bin/time
#./react-and-go
go run main.go
