package routers

import (
	"loan-management/api/controllers"
	"loan-management/api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sv-tools/mongoifc"
)

func AddUserRoutes(r *gin.Engine, db mongoifc.Database) {
	userController := controllers.NewUserController(db)
	userRouteGroup := r.Group("/users")
	{
		userRouteGroup.POST("/register", userController.SignUp)
		userRouteGroup.GET("/verify-email", userController.VerifyEmail)
		userRouteGroup.POST("/login", userController.Login)
		userRouteGroup.POST("/password-reset", userController.ForgetPassword)
		userRouteGroup.POST("/password-update", userController.ResetPassword)
		userRouteGroup.POST("/token/refresh", userController.RefreshAccessToken)
		userRouteGroup.GET("/profile", middlewares.JWTMiddleware(), userController.GetProfile)
	}
	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middlewares.JWTMiddleware())
	adminRoutes.Use(middlewares.AdminMiddleware())
	{
		adminRoutes.GET("/users", userController.GetAllUsers)
		adminRoutes.GET("/users/:id", userController.GetUserByID)
	}
}
