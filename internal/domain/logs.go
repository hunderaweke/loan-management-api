package domain

import "time"

const SystemLogCollection = "system_logs"

type SystemLog struct {
	ID        string    `bson:"_id" json:"id"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Category  string    `bson:"category" json:"category"`
	Message   string    `bson:"message" json:"message"`
}
type LogRepository interface {
	Create(log SystemLog) error
	GetAll() ([]SystemLog, error)
}
type LogUsecase interface {
	GetAll() ([]SystemLog, error)
}
