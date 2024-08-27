package controllers

import (
	"loan-management/internal/domain"
	"loan-management/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sv-tools/mongoifc"
)

type LoanController struct {
	loanUsecase domain.LoanUsecase
}

func NewLoanController(db mongoifc.Database) LoanController {
	usecase := usecases.NewLoanUsecase(db)
	return LoanController{loanUsecase: usecase}
}

func (c *LoanController) CreateLoan(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	loan := domain.Loan{}
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid data format"})
		return
	}
	loan, err := c.loanUsecase.CreateLoan(userID.(string), loan.Ammount)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, loan)
}

func (c *LoanController) ViewLoanStatus(ctx *gin.Context) {
	loanID := ctx.Param("id")
	loan, err := c.loanUsecase.ViewLoanStatus(loanID)
	userID, _ := ctx.Get("userID")
	if err != nil || loan.UserID != userID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}
	ctx.JSON(http.StatusOK, loan)
}

func (c *LoanController) ViewAllLoans(ctx *gin.Context) {
	status := ctx.Query("status")
	order := ctx.Query("order")
	filter := map[string]string{}
	isAdmin, _ := ctx.Get("isAdmin")
	if !isAdmin.(bool) {
		userID, _ := ctx.Get("userID")
		filter["user_id"] = userID.(string)
	}
	if status != "" {
		filter["status"] = status
		if order != "" {
			filter["order"] = order
		}
	}
	loans, err := c.loanUsecase.ViewAllLoans(filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, loans)
}

func (c *LoanController) ApproveRejectLoan(ctx *gin.Context) {
	loanID := ctx.Param("id")
	st := ctx.Param("status")
	if st != "approve" && st != "reject" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	if st == "approve" {
		st = "approved"
	} else {
		st = st + "d"
	}
	updatedLoan, err := c.loanUsecase.ApproveRejectLoan(loanID, st)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedLoan)
}

func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	loanID := ctx.Param("id")
	err := c.loanUsecase.DeleteLoan(loanID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "loan deleted successfully"})
}
