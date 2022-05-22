package config

func AccessiblePermissions() map[string]string {
	const connectionService = "/connection_service.ConnectionService/"

	return map[string]string{
		connectionService + "Register":          "RegisterConnection",
		connectionService + "AddConnection":     "CreateConnection",
		connectionService + "RejectConnection":  "RejectConnection",
		connectionService + "ApproveConnection": "ApproveConnection",
	}
}
