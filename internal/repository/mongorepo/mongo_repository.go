package mongorepo

import (
	"company-service/internal/domain"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

// Cria uma nova empresa na base de dados
func (r *mongoRepository) Create(ctx context.Context, company *domain.Company) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	company.BeforeCreate()

	result, err := r.collection.InsertOne(ctx, company)
	if err != nil {
		return err
	}

	// Define ID gerado pelo MongoDB
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		company.ID = oid.Hex()
	}
	return nil
}

// GetByID busca uma empresa pelo seu ID.
func (r *mongoRepository) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid company ID")
	}

	var company domain.Company
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// GetByCNPJ busca uam empresa pelo CNPJ.
func (r *mongoRepository) GetByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var company domain.Company
	err := r.collection.FindOne(ctx, bson.M{"cnpj": cnpj}).Decode(&company)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

// Update atualiza uma empresa existente.
func (r *mongoRepository) Update(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(company.ID)
	if err != nil {
		return nil, errors.New("invalid company ID")
	}

	company.BeforeUpdate()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"cnpj":                            company.CNPJ,
			"fantasy_name":                    company.FantasyName,
			"corporate_name":                  company.CorporateName,
			"address":                         company.Address,
			"employee_count":                  company.EmployeeCount,
			"required_min_pwd_employee_count": company.RequiredMinPWDEmployeeCount,
			"updated_at":                      company.UpdatedAt,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCompany domain.Company
	err = r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCompany)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("company not found")
		}
		return nil, err
	}

	return &updatedCompany, nil
}

// Delete remove uma empresa pelo ID.
func (r *mongoRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid company ID")
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("company not found")
	}

	return nil
}

// List lista empresas com paginação.
func (r *mongoRepository) List(ctx context.Context, page int, limit int) ([]*domain.Company, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetSort(bson.D{{Key: "created_at", Value: -1}}) // prioridade de exibição para os inseridos mais recentes

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var companies []*domain.Company
	for cursor.Next(ctx) {
		var company domain.Company
		if err := cursor.Decode(&company); err != nil {
			return nil, err
		}
		companies = append(companies, &company)
	}

	return companies, nil
}

// Count implements repository.CompanyRepository.
func (r *mongoRepository) Count(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{})
}
