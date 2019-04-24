#!/bin/sh
docker run -it --rm -v$(pwd):/work -w /work golang:alpine sh -c'
    apk add --update git zip
    GOOS=linux go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o main main.go
    zip main.zip main
'