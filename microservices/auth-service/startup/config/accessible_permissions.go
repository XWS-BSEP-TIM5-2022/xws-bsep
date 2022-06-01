package config

func AccessiblePermissions() map[string]string {
	const authService = "/auth_service.AuthService/"

	return map[string]string{
		authService + "UpdateUsername":    "UpdateUsername",
		authService + "ChangePassword":    "UpdatePassword",
		authService + "AdminsEndpoint":    "AdminsEndpoint", // empty endpoint
		authService + "CreateNewAPIToken": "CreateNewAPIToken",
		//authService + "GetUsernameByApiToken": "GetUsernameByApiToken",
	}
}
