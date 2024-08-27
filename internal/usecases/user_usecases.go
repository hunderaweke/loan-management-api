package usecases

import (
	"errors"
	"fmt"
	"loan-management/internal/domain"
	"loan-management/internal/repositories"
	"loan-management/pkg/infrastructures"
	"time"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(db mongoifc.Database) domain.UserUsecases {
	repo := repositories.NewUserRepository(db)
	return &userUsecase{userRepository: repo}
}

func (uc *userUsecase) Register(user domain.User) (domain.User, error) {
	hashedPassword, err := infrastructures.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword
	user.ID = primitive.NewObjectIDFromTimestamp(time.Now()).Hex()
	expirationTime := time.Now().Add(1 * time.Hour)
	verificationToken, err := infrastructures.GenerateVerificationToken(user.ID, user.Email, "emailVerification", expirationTime)
	if err != nil {
		return domain.User{}, fmt.Errorf("generating token:", err)
	}
	go infrastructures.SendVerificationEmail(user.Email, verificationToken)
	return uc.userRepository.Create(user)
}

func (uc *userUsecase) VerifyEmail(token, email string) error {
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	if err := infrastructures.ValidateVerificationToken(token, email, user.ID, "emailVerification"); err != nil {
		return err
	}
	uc.userRepository.Update(user.ID, domain.User{IsActive: true})
	return nil
}

func (uc *userUsecase) Login(email, password string) (domain.User, error) {
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return user, err
	}
	if !user.IsActive {
		return domain.User{}, fmt.Errorf("user account is not activated")
	}
	if infrastructures.ComparePassword(user.Password, password) {
		return user, nil
	}
	return domain.User{}, fmt.Errorf("incorrect password or email")
}

func (uc *userUsecase) GetProfile(userID string) (domain.User, error) {
	return uc.userRepository.GetByID(userID)
}

func (uc *userUsecase) ForgetPassword(email string) error {
	fmt.Println(email)
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return errors.New("reseting unactivated account is not allowed")
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := infrastructures.GenerateVerificationToken(user.ID, email, "ForgetPassword", expirationTime)
	go infrastructures.SendPasswordResetEmail(email, token)
	return nil
}

func (uc *userUsecase) ResetPassword(token, email, newPassword string) error {
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	if err := infrastructures.ValidateVerificationToken(token, email, user.ID, "ForgetPassword"); err != nil {
		return err
	}
	hashedPassword, err := infrastructures.HashPassword(newPassword)
	if err != nil {
		return err
	}
	uc.userRepository.Update(user.ID, domain.User{Password: hashedPassword})
	return nil
}

func (uc *userUsecase) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := infrastructures.ValidateJWTToken(refreshToken)
	if err != nil {
		return "", err
	}
	user, err := uc.userRepository.GetByID(claims.UserID)
	if err != nil {
		return "", err
	}
	accessToken, err := infrastructures.GenerateJWTToken(user, 1*time.Hour)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (uc *userUsecase) GetAllUsers() ([]domain.User, error) {
	return uc.userRepository.Get()
}

func (uc *userUsecase) GetUserByID(userID string) (domain.User, error) {
	return uc.userRepository.GetByID(userID)
}
