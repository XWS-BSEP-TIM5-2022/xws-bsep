package config

func AccessiblePermissions() map[string]string {
	const postService = "/post_service.PostService/"

	return map[string]string{
		postService + "Get":         "GetPostById",
		postService + "GetAll":      "GetAllPosts",
		postService + "Update":      "UpdatePost",
		postService + "Insert":      "CreatePost",
		postService + "LikePost":    "UpdatePostLikes",
		postService + "DislikePost": "UpdatePostDislikes",
		postService + "CommentPost": "UpdatePostComments",
		postService + "NeutralPost": "NeutralPost",
	}
}
