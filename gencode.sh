#!/bin/sh

protoc --proto_path=proto --go-grpc_out=./ model.proto 
protoc --proto_path=proto --go-_out=./ model.proto
