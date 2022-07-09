package config

// TODO SD: obrisati
func AccessibleRoles() map[string][]string {
	const userService = "/user_service.UserService/"

	return map[string][]string{
		//userService + "Get":                          {"User"},	// TODO: TM
		userService + "GetAll": {"User"},
		//userService + "Update":                       {"User"},
		userService + "UpdateBasicInfo":              {"User"},
		userService + "UpdatePostNotification":       {"User"},
		userService + "UpdateExperienceAndEducation": {"User"},
		userService + "UpdateSkillsAndInterests":     {"User"},
		userService + "GetLoggedInUserInfo":          {"User"},
		userService + "UpdatePrivacy":                {"User"},
	}
}
