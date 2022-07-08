package config

func AccessiblePermissions() map[string]string {
	const messageService = "/message_service.MessageService/"

	return map[string]string{
		messageService + "GetConversation":            "CreateMessage",
		messageService + "GetAllConversationsForUser": "CreateMessage",
		messageService + "NewMessage":                 "CreateMessage",
	}
}
