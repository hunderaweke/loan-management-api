package usecases

import (
	"loan-management/internal/domain"
	"loan-management/internal/repositories"

	"github.com/sv-tools/mongoifc"
)

type logUsecase struct {
	logRepository domain.LogRepository
}

func NewLogUsecase(db mongoifc.Database) domain.LogUsecase {
	logRepo := repositories.NewLogRepository(db)
	return &logUsecase{logRepository: logRepo}
}

func (uc *logUsecase) GetAll() ([]domain.SystemLog, error) {
	return uc.logRepository.GetAll()
}
