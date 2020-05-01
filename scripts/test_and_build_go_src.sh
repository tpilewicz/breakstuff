#!/bin/bash

cd ./components/clausius/src
go test ./...
if [[ $? != 0 ]]; then
    echo "Tests failed. Won't build!"
    exit 1
fi

for filename in get_grid set_cell troublemaker; do
    echo "Building $filename"
    cd ./$filename
    GOOS=linux go build main.go
    cd ..
done
