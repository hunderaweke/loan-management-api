package repositories

import (
	"context"
	"loan-management/internal/domain"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
)

type logRepository struct {
	collection mongoifc.Collection
}

func NewLogRepository(db mongoifc.Database) domain.LogRepository {
	collection := db.Collection("system_logs")
	return &logRepository{collection: collection}
}

func (r *logRepository) Create(log domain.SystemLog) error {
	_, err := r.collection.InsertOne(context.TODO(), log)
	return err
}

func (r *logRepository) GetAll() ([]domain.SystemLog, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var logs []domain.SystemLog
	if err := cursor.All(context.TODO(), &logs); err != nil {
		return nil, err
	}

	return logs, nil
}
