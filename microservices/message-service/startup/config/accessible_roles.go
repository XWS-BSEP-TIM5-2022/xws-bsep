package config

func AccessibleRoles() map[string][]string {
	const messageService = "/message_service.MessageService/"

	return map[string][]string{
		messageService + "GetConversation":            {"User"},
		messageService + "GetAllConversationsForUser": {"User"},
		messageService + "NewMessage":                 {"User"},
	}
}
