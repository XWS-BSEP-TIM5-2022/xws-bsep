#!/bin/(shell)

protoc -I ./auth_service \
       --go_out ./auth_service --go_opt paths=source_relative \
       --go-grpc_out ./auth_service --go-grpc_opt paths=source_relative \
	--grpc-gateway_out ./auth_service --grpc-gateway_opt paths=source_relative \
       ./auth_service/auth_service.proto
