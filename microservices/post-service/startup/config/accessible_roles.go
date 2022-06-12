package config

// => /naziv_paketa.naziv_servisa/naziv_metode : { naziv_role1, naziv_role2 }
// za sve metode koje ne treba da se presrecu -> ne dodaju se u mapu
func AccessibleRoles() map[string][]string {
	const postService = "/post_service.PostService/"

	return map[string][]string{
		//postService + "Get": {"User"},			// TODO: TM
		//postService + "GetAll":      {"User"}, 	// TODO: TM
		postService + "Update":      {"User"},
		postService + "Insert":      {"User"},
		postService + "LikePost":    {"User"},
		postService + "DislikePost": {"User"},
		postService + "CommentPost": {"User"},
		postService + "NeutralPost": {"User"},
		//postService + "GetAllByUser": {"User"},
	}
}
