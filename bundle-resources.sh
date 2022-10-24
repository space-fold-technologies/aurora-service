#!/usr/bin/env bash

go-bindata -o ./app/core/server/resources.go ./resources/ ./resources/migrations
sed 's/package main/package server/g' ./app/core/server/resources.go > ./app/core/server/content.go
rm ./app/core/server/resources.go
mv ./app/core/server/content.go ./app/core/server/resources.go