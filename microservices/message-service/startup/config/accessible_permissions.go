package config

func AccessiblePermissions() map[string]string {
	const messageService = "/message_service.MessageService/"

	return map[string]string{
		//messageService + "GetAll": "GetAllUsers",
	}
}
