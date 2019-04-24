package main

import (
	"github.com/Xinguang/lambda-sendmail/recaptcha"
	"github.com/aws/aws-lambda-go/lambda"
)

// Golang Alpine images doesn't build statically linked binary:
// ELF 64-bit LSB executable, x86-64, version 1 (SYSV), dynamically linked, interpreter /lib/ld-musl-x86_64.so.1, stripped

// The interpreter doesn't exist in Lambda environment so that's why it throws "no such file or directory error".

// To solve that modify your Makefile to use
// env GOOS=linux go build -v -ldflags '-d -s -w' -a -tags netgo -installsuffix netgo -o main main.go as a build command.
func main() {
	// Make the handler available for Remote Procedure Call by AWS Lambda
	lambda.Start(recaptcha.Handle)
}
