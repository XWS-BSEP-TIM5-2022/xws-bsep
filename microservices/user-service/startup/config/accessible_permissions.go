package config

func AccessiblePermissions() map[string]string {
	const userService = "/user_service.UserService/"

	return map[string]string{
		//userService + "Get":                          "GetUserById",	// TODO: TM
		userService + "GetAll": "GetAllUsers",
		//userService + "Update":                       "UpdateUserProfile",
		userService + "UpdatePostNotification":       "UpdateUserProfile",
		userService + "UpdateBasicInfo":              "UpdateUserProfile",
		userService + "UpdateExperienceAndEducation": "UpdateUserProfile",
		userService + "UpdateSkillsAndInterests":     "UpdateUserProfile",
		userService + "GetLoggedInUserInfo":          "GetLoggedInUserInfo",
		userService + "UpdatePrivacy":                "UpdateUserProfile",

		//userService + "GetIdByUsername":              "GetIdByUsername",
	}
}
