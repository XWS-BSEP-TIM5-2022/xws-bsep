package config

func AccessibleRoles() map[string][]string {
	const eventService = "/event_service.EventService/"

	return map[string][]string{
		eventService + "GetAllEvents": {"Admin"}}
}
