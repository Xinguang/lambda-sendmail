#!/bin/sh
# use sam
# docker run -it --rm -v $(pwd):/work -w /work -v $(pwd)/.aws:/root/.aws lambci/lambda:build-go1.x bash
# docker run -it --rm -v $(pwd):/work -w /work lambci/lambda:build-go1.x sam build

docker run -it --rm -v$(pwd):/work -w /work golang:alpine sh -c'
    apk add --update git zip
    GOOS=linux go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o main main.go
    zip main.zip main
'