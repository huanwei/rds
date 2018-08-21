#!/bin/bash

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -i -o ../../docker/mysql-operator/mysql-operator  main.go
