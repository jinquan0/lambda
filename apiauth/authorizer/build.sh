#!/bin/bash
 
echo "Build golang binary"
go mod init cmd/binary/main
go mod tidy
GOOS=linux GOARCH=amd64 go build -o main cmd/*.go
 
echo "Create a ZIP file"
zip deployment.zip main
 
echo "Cleaning up"
rm main go.mod go.sum