package domain

import "time"

const LoanColletion = "loans"

type Loan struct {
	ID        string
	UserID    string
	Ammount   string
	CreatedAt time.Time
	Status    string
}

type LoanRepository interface {
	GetByID(id string) (Loan, error)
	Get(filter map[string]string) ([]Loan, error)
	Delete(loanID string) error
	Update(id string, updateData Loan) (Loan, error)
	Create(loan Loan) (Loan, error)
}

type LoanUsecase interface {
	ApplyForLoan(userID string, amount string) (Loan, error)
	ViewLoanStatus(id string) (Loan, error)
	ViewAllLoans(filter map[string]string) ([]Loan, error)
	ApproveRejectLoan(id string, status string) (Loan, error)
	DeleteLoan(id string) error
}
