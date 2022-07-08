package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
)

func mapNotification(notification *domain.Notification) *pb.Notification {
	notificationPb := &pb.Notification{
		Id:   notification.Id.Hex(),
		Time: notification.Time.String(),
	}

	return notificationPb
}

func removeMalicious(value string) string {

	var lenId = len(value)
	var checkId = ""
	for i := 0; i < lenId; i++ {
		char := string(value[i])
		if char != "$" {
			checkId = checkId + char
		}
	}
	return checkId
}
