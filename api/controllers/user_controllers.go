package controllers

import (
	"loan-management/internal/domain"
	"loan-management/internal/usecases"
	"loan-management/pkg/infrastructures"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sv-tools/mongoifc"
)

type UserController struct {
	userUsecase domain.UserUsecases
}

func NewUserController(db mongoifc.Database) UserController {
	usecase := usecases.NewUserUsecase(db)
	return UserController{userUsecase: usecase}
}

func (uc *UserController) SignUp(ctx *gin.Context) {
	var user domain.User

	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "invalid data format"})
		return
	}
	_, err := uc.userUsecase.Register(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "verification email sent to your email address"})
}

func (uc *UserController) VerifyEmail(ctx *gin.Context) {
	token := ctx.Query("token")
	email := ctx.Query("email")
	if err := uc.userUsecase.VerifyEmail(token, email); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "account activated successfully"})
}

func (uc *UserController) Login(ctx *gin.Context) {
	credentials := struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "invalid credentials"})
		return
	}
	user, err := uc.userUsecase.Login(credentials.Email, credentials.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	refreshToken, err := infrastructures.GenerateJWTToken(user, (7 * time.Hour))
	accessToken, err := infrastructures.GenerateJWTToken(user, (1 * time.Hour))
	ctx.JSON(http.StatusOK, gin.H{"refresh_token": refreshToken, "access_token": accessToken})
}

func (uc *UserController) RefreshAccessToken(ctx *gin.Context) {
	token := struct {
		RefreshToken string `json:"refresh_token"`
	}{}
	if err := ctx.ShouldBind(&token); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "refresh_token field is required"})
		return
	}
	accessToken, err := uc.userUsecase.RefreshAccessToken(token.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (uc *UserController) ForgetPassword(ctx *gin.Context) {
	email := ctx.Query("email")
	err := uc.userUsecase.ForgetPassword(string(email))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "password reset link sent to your email"})
}

func (uc *UserController) ResetPassword(ctx *gin.Context) {
	email := ctx.Query("email")
	token := ctx.Query("token")
	newPassword := struct {
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}{}
	if err := ctx.ShouldBind(&newPassword); err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "unacceptable data format"})
		return
	}
	if newPassword.Password != newPassword.ConfirmPassword {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"error": "password and confirm password don match"})
		return
	}
	err := uc.userUsecase.ResetPassword(token, email, newPassword.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusAccepted, gin.H{"message": "password update successful"})
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	user, err := uc.userUsecase.GetProfile(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("id")
	user, err := uc.userUsecase.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := uc.userUsecase.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
