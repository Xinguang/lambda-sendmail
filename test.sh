#!/bin/sh
# test
# docker run -it --rm -v $(pwd):$(pwd) -w $(pwd) -v $(pwd)/.aws:/root/.aws -v /var/run/docker.sock:/var/run/docker.sock lambci/lambda:build-go1.x bash
# build
docker run -it --rm -v $(pwd):/work -w /work lambci/lambda:build-go1.x sam build 
# invoke
sam local invoke SendmailFunction --event test-event.json --env-vars .env.test