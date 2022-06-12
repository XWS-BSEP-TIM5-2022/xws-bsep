package config

func AccessiblePermissions() map[string]string {
	const postService = "/post_service.PostService/"

	return map[string]string{
		//postService + "Get": "GetPostById",			// TODO: TM
		//postService + "GetAll":      "GetAllPosts",	// TODO: TM
		postService + "Update":      "UpdatePost",
		postService + "Insert":      "CreatePost",
		postService + "LikePost":    "UpdatePostLikes",
		postService + "DislikePost": "UpdatePostDislikes",
		postService + "CommentPost": "UpdatePostComments",
		postService + "NeutralPost": "NeutralPost",
		//postService + "InsertJobOffer": "InsertJobOffer",
	}
}
