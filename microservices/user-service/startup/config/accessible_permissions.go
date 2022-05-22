package config

func AccessiblePermissions() map[string]string {
	const userService = "/user_service.UserService/"

	return map[string]string{
		userService + "Get":                          "GetUserById",
		userService + "GetAll":                       "GetAll",
		userService + "Update":                       "UpdateUserProfile",
		userService + "UpdateBasicInfo":              "UpdateUserProfile",
		userService + "UpdateExperienceAndEducation": "UpdateUserProfile",
		userService + "UpdateSkillsAndInterests":     "UpdateUserProfile",
		userService + "GetLoggedInUserInfo":          "GetLoggedInUserInfo",
	}
}