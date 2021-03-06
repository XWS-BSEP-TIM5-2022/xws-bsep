#!/bin/(shell)
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

protoc -I ./user_service \
       --go_out ./user_service --go_opt paths=source_relative \
       --go-grpc_out ./user_service --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./user_service --grpc-gateway_opt paths=source_relative \
       ./user_service/user_service.proto


protoc -I ./auth_service \
       --go_out ./auth_service --go_opt paths=source_relative \
       --go-grpc_out ./auth_service --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./auth_service --grpc-gateway_opt paths=source_relative \
       ./auth_service/auth_service.proto


