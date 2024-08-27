package controllers

import (
	"loan-management/internal/domain"
	"loan-management/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sv-tools/mongoifc"
)

type LogController struct {
	logUsecase domain.LogUsecase
}

func NewLogController(db mongoifc.Database) *LogController {
	usecase := usecases.NewLogUsecase(db)
	return &LogController{logUsecase: usecase}
}

func (c *LogController) GetLogs(ctx *gin.Context) {
	logs, err := c.logUsecase.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, logs)
}
