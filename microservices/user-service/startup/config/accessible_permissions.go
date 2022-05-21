package config

func AccessiblePermissions() map[string]string {
	const userService = "/user_service.UserService/"

	return map[string]string{
		userService + "Get":                          "Get",
		userService + "GetAll":                       "GetAll",
		userService + "Update":                       "Update",
		userService + "UpdateBasicInfo":              "UpdateUserProfile",
		userService + "UpdateExperienceAndEducation": "UpdateUserProfile",
		userService + "UpdateSkillsAndInterests":     "UpdateUserProfile",
		userService + "GetLoggedInUserInfo":          "GetLoggedInUserInfo",
	}
}
