module github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/api-gateway

go 1.17

replace github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common => ../common

require (
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.0
	github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common v0.1.0
	google.golang.org/grpc v1.46.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220111092808-5a964db01320 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220429170224-98d788798c3e // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
