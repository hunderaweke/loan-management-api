package repositories

import (
	"context"
	"fmt"
	"loan-management/internal/domain"

	"github.com/sv-tools/mongoifc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type loanRepository struct {
	collection mongoifc.Collection
}

func NewLoanRepository(db mongoifc.Database) domain.LoanRepository {
	collection := db.Collection(domain.LoanColletion)
	return &loanRepository{collection: collection}
}

func (r *loanRepository) GetByID(id string) (domain.Loan, error) {
	var loan domain.Loan
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(context.TODO(), filter).Decode(&loan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Loan{}, ErrUserNotFound
		}
		return domain.Loan{}, err
	}
	return loan, nil
}

func (r *loanRepository) Get(filter map[string]string) ([]domain.Loan, error) {
	filterOptions := bson.M{}
	userID, ok := filter["user_id"]
	if ok {
		filterOptions["user_id"] = userID
	}
	status, ok := filter["status"]
	if ok {
		filterOptions["status"] = status
	}
	order, ok := filter["order"]
	if !ok {
		if status == "pending" {
			order = "asc"
		} else {
			order = "desc"
		}
	}
	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}
	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: sortOrder}})
	cursor, err := r.collection.Find(context.TODO(), filterOptions, findOptions)
	if err != nil {
		return []domain.Loan{}, err
	}
	defer cursor.Close(context.TODO())

	var loans []domain.Loan
	for cursor.Next(context.TODO()) {
		var l domain.Loan
		if err := cursor.Decode(&l); err != nil {
			return loans, nil
		}
		loans = append(loans, l)
	}
	return loans, nil
}

func (r *loanRepository) Delete(loanID string) error {
	filter := bson.M{"_id": loanID}

	result, err := r.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("loan not found")
	}

	return nil
}

func (r *loanRepository) Update(id string, updateData domain.Loan) (domain.Loan, error) {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{}}
	if updateData.Status != "" {
		update["$set"].(bson.M)["status"] = updateData.Status
	}
	if updateData.Ammount != "" {
		update["$set"].(bson.M)["amount"] = updateData.Ammount
	}
	if updateData.Status != "" {
		update["$set"].(bson.M)["status"] = updateData.Status
	}
	_, err := r.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("failed to update loan: %v", err)
	}

	updatedLoan, err := r.GetByID(id)
	if err != nil {
		return domain.Loan{}, fmt.Errorf("failed to retrieve updated loan: %v", err)
	}

	return updatedLoan, nil
}

func (r *loanRepository) Create(loan domain.Loan) (domain.Loan, error) {
	loan.ID = primitive.NewObjectID().Hex()

	_, err := r.collection.InsertOne(context.TODO(), loan)
	if err != nil {
		return domain.Loan{}, err
	}

	return loan, nil
}
