package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/notification_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/notification_service/domain"
	"time"
)

func mapNotification(notification *domain.Notification) *pb.Notification {
	notificationPb := &pb.Notification{
		Id:     notification.Id.Hex(),
		Date:   notification.Date.String(),
		Text:   notification.Text,
		UserId: notification.UserId,
		Type:   pb.Notification_NotificationTypeEnum(notification.Type), //TODO: proveriti!
		Read:   notification.Read,
	}

	return notificationPb
}

func mapInsertNotification(notification *pb.Notification) (*domain.Notification, error) {
	postPb := &domain.Notification{
		Date:   time.Now(),
		Text:   notification.Text,
		UserId: notification.UserId,
		Read:   notification.Read,
		Type:   mapInsertNotificationType(notification.Type),
	}

	return postPb, nil
}

func mapInsertNotificationType(notif_type pb.Notification_NotificationTypeEnum) domain.NotificationTypeEnum {
	switch notif_type {
	case pb.Notification_Follow:
		return domain.Follow
	case pb.Notification_Message:
		return domain.Message
	case pb.Notification_Post:
		return domain.Post
	}
	return -1
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
