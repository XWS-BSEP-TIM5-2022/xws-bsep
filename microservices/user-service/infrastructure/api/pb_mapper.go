package api

import "user-service/domain"

func mapUser(order *domain.User) *pb.User {
	userPb := &pb.User{
		Id:   order.Id.Hex(),
		Name: order.Name,
	}
	return userPb
}
