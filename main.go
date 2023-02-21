package main

import (
	"fmt"
	"time"

	"github.com/aldisaputra17/posting-board/controllers"
	"github.com/aldisaputra17/posting-board/databases"

	// "github.com/aldisaputra17/posting-board/dummy"
	"github.com/aldisaputra17/posting-board/repositories"
	"github.com/aldisaputra17/posting-board/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	contextTimeout    time.Duration                  = 10 * time.Second
	db                *gorm.DB                       = databases.ConnectDB()
	userRepository    repositories.UserRepository    = repositories.NewUserRepository(db)
	postRepository    repositories.PostRepository    = repositories.NewPostRepository(db)
	commentRepository repositories.CommentRepository = repositories.NewCommentRepository(db)
	likeRepository    repositories.LikeRepository    = repositories.NewLikeRepository(db)
	authService       services.AuthService           = services.NewAuthService(userRepository, contextTimeout)
	jwtService        services.JWTService            = services.NewJWTService()
	postService       services.PostService           = services.NewPostService(postRepository, contextTimeout)
	commentService    services.CommentService        = services.NewCommentService(commentRepository, contextTimeout)
	likeService       services.LikeService           = services.NewLikeService(likeRepository, contextTimeout)
	authController    controllers.AuthController     = controllers.NewAuthController(authService, jwtService)
	postController    controllers.PostController     = controllers.NewPostController(postService, jwtService)
	commentController controllers.CommentController  = controllers.NewCommentController(commentService, jwtService)
	likeController    controllers.LikeController     = controllers.NewLikeController(likeService, jwtService)
)

func main() {
	fmt.Println("Starting Server")
	defer databases.CloseDatabaseConnection(db)
	// dummy.CreateCrypto(db)

	r := gin.Default()

	api := r.Group("api")

	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}
	postRoutes := api.Group("/post")
	{
		postRoutes.POST("", postController.Add)
		postRoutes.GET("", postController.List)
		postRoutes.GET("/public", postController.ListPublic)
		postRoutes.DELETE("/:id", postController.Deleted)
		postRoutes.GET("/comment", postController.CommentMost)
		postRoutes.GET("/like", postController.LikeMost)
		postRoutes.GET("/user/:user_id", postController.FindByUserID)
	}
	commentRoutes := api.Group("/comment")
	{
		commentRoutes.POST("", commentController.Add)
		commentRoutes.GET("", commentController.List)
		commentRoutes.GET("/:post_id", commentController.ListUserComment)
	}
	likeRoutes := api.Group("/like")
	{
		likeRoutes.POST("", likeController.Add)
		likeRoutes.GET("/:post_id", likeController.ListUser)
	}
	r.Run()
}
