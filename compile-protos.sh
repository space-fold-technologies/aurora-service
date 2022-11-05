#!/usr/bin/env bash

protoc -I protos/ protos/apps.proto --go_out=app/domain/apps
protoc -I protos/ protos/clusters.proto --go_out=app/domain/clusters
protoc -I protos/ protos/nodes.proto --go_out=app/domain/nodes
protoc -I protos/ protos/envs.proto --go_out=app/domain/environments
protoc -I protos/ protos/users.proto --go_out=app/domain/users
protoc -I protos/ protos/teams.proto --go_out=app/domain/teams
protoc -I protos/ protos/authorization.proto --go_out=app/domain/authorization