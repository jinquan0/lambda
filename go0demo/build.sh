#!/bin/bash
 
echo "Build the binary"
go mod init cmd/main
go mod tidy
GOOS=linux GOARCH=amd64 go build -o main cmd/main.go
 
echo "Create a ZIP file"
zip deployment.zip main
 
echo "Cleaning up"
rm main
