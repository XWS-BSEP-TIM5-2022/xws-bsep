package api

import (
	pb "github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/common/proto/event_service"
	"github.com/XWS-BSEP-TIM5-2022/xws-bsep/microservices/event_service/domain"
)

func mapEvent(event *domain.Event) *pb.Event {
	eventPb := &pb.Event{
		Id:   event.Id.Hex(),
		Date: event.Date.String(),
	}

	return eventPb
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
