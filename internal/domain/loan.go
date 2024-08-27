package domain

import "time"

const LoanColletion = "loans"

type Loan struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Ammount   string    `json:"ammount" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	Status    string    `json:"status"`
}

type LoanRepository interface {
	GetByID(id string) (Loan, error)
	Get(filter map[string]string) ([]Loan, error)
	Delete(loanID string) error
	Update(id string, updateData Loan) (Loan, error)
	Create(loan Loan) (Loan, error)
}

type LoanUsecase interface {
	CreateLoan(userID string, amount string) (Loan, error)
	ViewLoanStatus(id string) (Loan, error)
	ViewAllLoans(filter map[string]string) ([]Loan, error)
	ApproveRejectLoan(id string, status string) (Loan, error)
	DeleteLoan(id string) error
}
