#!/bin/sh

for GOARCH in 386 amd64; do
	GOOS=linux GOARCH=$GOARCH go build -o build/dbignore-linux-$GOARCH ./dbignore.go
done
