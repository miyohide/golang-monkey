#!/bin/bash
 
APP_DIR="/go/src/github.com/miyohide/monkey/"
 
mkdir -p ${APP_DIR} && cp -r ./ ${APP_DIR} && cd ${APP_DIR}
 
echo "#######################"
echo "# Running Test"
go test ./...
