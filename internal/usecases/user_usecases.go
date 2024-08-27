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
	logRepository  domain.LogRepository
}

func NewUserUsecase(db mongoifc.Database) domain.UserUsecases {
	userRepo := repositories.NewUserRepository(db)
	logRepo := repositories.NewLogRepository(db)
	return &userUsecase{
		userRepository: userRepo,
		logRepository:  logRepo,
	}
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
		return domain.User{}, fmt.Errorf("generating token: %v", err)
	}

	go infrastructures.SendVerificationEmail(user.Email, verificationToken)

	createdUser, err := uc.userRepository.Create(user)
	if err != nil {
		return domain.User{}, err
	}

	// Log the registration
	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "User Registration",
		Message:   fmt.Sprintf("User %s registered with email %s", user.ID, user.Email),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log user registration: %v\n", err)
	}

	return createdUser, nil
}

// VerifyEmail verifies the user's email with the token
func (uc *userUsecase) VerifyEmail(token, email string) error {
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	if err := infrastructures.ValidateVerificationToken(token, email, user.ID, "emailVerification"); err != nil {
		return err
	}
	if _, err := uc.userRepository.Update(user.ID, domain.User{IsActive: true}); err != nil {
		return err
	}
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
	if !infrastructures.ComparePassword(user.Password, password) {
		return domain.User{}, fmt.Errorf("incorrect password or email")
	}

	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Login Attempt",
		Message:   fmt.Sprintf("User %s logged in successfully", user.ID),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log successful login attempt: %v\n", err)
	}

	return user, nil
}

func (uc *userUsecase) GetProfile(userID string) (domain.User, error) {
	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (uc *userUsecase) ForgetPassword(email string) error {
	user, err := uc.userRepository.GetByEmail(email)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return errors.New("resetting unactivated account is not allowed")
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := infrastructures.GenerateVerificationToken(user.ID, email, "ForgetPassword", expirationTime)
	if err != nil {
		return err
	}
	go infrastructures.SendPasswordResetEmail(email, token)
	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Password Reset Request",
		Message:   fmt.Sprintf("User %s requested a password reset", user.ID),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log password reset request: %v\n", err)
	}

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
	if _, err := uc.userRepository.Update(user.ID, domain.User{Password: hashedPassword}); err != nil {
		return err
	}

	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Password Reset Completion",
		Message:   fmt.Sprintf("User %s completed password reset", user.ID),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log password reset completion: %v\n", err)
	}

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
	users, err := uc.userRepository.Get()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (uc *userUsecase) GetUserByID(userID string) (domain.User, error) {
	user, err := uc.userRepository.GetByID(userID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}
