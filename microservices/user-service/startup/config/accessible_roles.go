package config

// => /naziv_modula_proto_file.naziv_servisa/naziv_metode : { naziv_role1, naziv_role2 }
// za sve metode koje ne treba da se presrecu -> ne dodaju se u mapu
func AccessibleRoles() map[string][]string {
	const userService = "/user.UserService/"
	// const postService = "/post.PostService/"

	return map[string][]string{
		userService + "Get":    {"User"},
		userService + "GetAll": {"User"},
		userService + "Update": {"User"},
	}
}
