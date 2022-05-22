package config

// TODO SD: delete

func AccessibleRoles() map[string][]string {
	const authService = "/auth_service.AuthService/"

	return map[string][]string{
		authService + "UpdateUsername": {"User"},
		authService + "ChangePassword": {"User"},
	}
}
