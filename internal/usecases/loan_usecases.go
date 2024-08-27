package usecases

import (
	"fmt"
	"loan-management/internal/domain"
	"loan-management/internal/repositories"
	"time"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type loanUsecase struct {
	loanRepository domain.LoanRepository
	logRepository  domain.LogRepository
}

func NewLoanUsecase(db mongoifc.Database) domain.LoanUsecase {
	loanRepo := repositories.NewLoanRepository(db)
	logRepo := repositories.NewLogRepository(db)
	return &loanUsecase{
		loanRepository: loanRepo,
		logRepository:  logRepo,
	}
}

func (uc *loanUsecase) DeleteLoan(id string) error {
	err := uc.loanRepository.Delete(id)
	if err != nil {
		return err
	}

	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Loan Deletion",
		Message:   fmt.Sprintf("Loan %s was deleted", id),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log loan deletion: %v\n", err)
	}

	return nil
}

// ApproveRejectLoan updates the status of a loan
func (uc *loanUsecase) ApproveRejectLoan(id string, status string) (domain.Loan, error) {
	loan, err := uc.loanRepository.GetByID(id)
	if err != nil {
		return domain.Loan{}, err
	}

	if loan.Status != "pending" {
		return domain.Loan{}, fmt.Errorf("loan status cannot be updated from %s", loan.Status)
	}

	loan.Status = status
	updatedLoan, err := uc.loanRepository.Update(id, loan)
	if err != nil {
		return domain.Loan{}, err
	}

	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Loan Application Status Update",
		Message:   fmt.Sprintf("Loan application %s was %s", id, status),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log loan status update: %v\n", err)
	}

	return updatedLoan, nil
}

// CreateLoan creates a new loan application
func (uc *loanUsecase) CreateLoan(userID string, amount string) (domain.Loan, error) {
	loan := domain.Loan{
		UserID:  userID,
		Ammount: amount,
		Status:  "pending",
	}
	createdLoan, err := uc.loanRepository.Create(loan)
	if err != nil {
		return domain.Loan{}, err
	}

	// Log the creation
	log := domain.SystemLog{
		ID:        primitive.NewObjectID().Hex(),
		Timestamp: time.Now(),
		Category:  "Loan Application Submission",
		Message:   fmt.Sprintf("User %s submitted a loan application with amount %s", userID, amount),
	}
	if err := uc.logRepository.Create(log); err != nil {
		fmt.Printf("Failed to log loan creation: %v\n", err)
	}

	return createdLoan, nil
}

// ViewAllLoans retrieves all loans based on filters
func (uc *loanUsecase) ViewAllLoans(filter map[string]string) ([]domain.Loan, error) {
	loans, err := uc.loanRepository.Get(filter)
	if err != nil {
		return nil, err
	}

	return loans, nil
}

// ViewLoanStatus retrieves a loan's status by ID
func (uc *loanUsecase) ViewLoanStatus(id string) (domain.Loan, error) {
	return uc.loanRepository.GetByID(id)
}
