package routers

import (
	"loan-management/api/controllers"
	"loan-management/api/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sv-tools/mongoifc"
)

func AddLoanRoutes(r *gin.Engine, db mongoifc.Database) {
	loanController := controllers.NewLoanController(db)
	loanRouter := r.Group("/loans")
	loanRouter.Use(middlewares.JWTMiddleware())
	{
		loanRouter.POST("/", loanController.CreateLoan)
		loanRouter.GET("/:id", loanController.ViewLoanStatus)
	}
	adminRouter := r.Group("/admin/loans")
	adminRouter.Use(middlewares.JWTMiddleware())
	adminRouter.Use(middlewares.AdminMiddleware())
	{
		adminRouter.GET("/", loanController.ViewAllLoans)
		adminRouter.PATCH("/:id/:status", loanController.ApproveRejectLoan)
		adminRouter.DELETE("/:id", loanController.DeleteLoan)
	}
}
