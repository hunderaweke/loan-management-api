package usecases

import (
	"fmt"
	"loan-management/internal/domain"
	"loan-management/internal/repositories"

	"github.com/sv-tools/mongoifc"
)

type loanUsecase struct {
	loanRepository domain.LoanRepository
}

func NewLoanRepository(db mongoifc.Database) domain.LoanUsecase {
	repo := repositories.NewLoanRepository(db)
	return &loanUsecase{loanRepository: repo}
}

func (uc *loanUsecase) DeleteLoan(id string) error {
	return uc.loanRepository.Delete(id)
}

func (uc *loanUsecase) ApproveRejectLoan(id string, status string) (domain.Loan, error) {
	loan, err := uc.loanRepository.GetByID(id)
	if err != nil {
		return domain.Loan{}, err
	}
	if loan.Status == "pending" {
		loan.Status = status
		return uc.loanRepository.Update(id, loan)
	}
	return domain.Loan{}, fmt.Errorf("loan cannot be updated")
}

func (uc *loanUsecase) ApplyForLoan(userID string, amount string) (domain.Loan, error) {
	loan := domain.Loan{
		UserID:  userID,
		Ammount: amount,
		Status:  "pending",
	}
	return uc.loanRepository.Create(loan)
}

func (uc *loanUsecase) ViewAllLoans(filter map[string]string) ([]domain.Loan, error) {
	return uc.loanRepository.Get(filter)
}

func (uc *loanUsecase) ViewLoanStatus(id string) (domain.Loan, error) {
	return uc.loanRepository.GetByID(id)
}
