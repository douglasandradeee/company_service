package mongorepo

import (
	"company-service/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func NewCompanyRepository(db *mongo.Database, collectionName string, timeout time.Duration) repository.CompanyRepository {
	return &mongoRepository{collection: db.Collection(collectionName), timeout: timeout}
}

func NewCompanyRepositoryWithTimeout(db *mongo.Database, collectionName string, timeout time.Duration) repository.CompanyRepository {
	return &mongoRepository{collection: db.Collection(collectionName), timeout: timeout}
}
