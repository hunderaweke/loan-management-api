package routers

import (
	"loan-management/api/controllers"
	"loan-management/api/middlewares"
	"loan-management/config"
	"loan-management/database"
	"log"

	"github.com/gin-gonic/gin"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.NewMongoDatabase(config)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()
	logController := controllers.NewLogController(db)
	router.GET("/admin/logs", middlewares.JWTMiddleware(), middlewares.AdminMiddleware(), logController.GetLogs)
	AddUserRoutes(router, db)
	AddLoanRoutes(router, db)
	router.Run(config.Server.Port)
}
