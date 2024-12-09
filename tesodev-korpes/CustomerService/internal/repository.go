package internal

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/pkg"
	"tesodev-korpes/shared/helpers"
	"tesodev-korpes/shared/model"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(col *mongo.Collection) *Repository {
	return &Repository{
		collection: col,
	}
}

func (r *Repository) FindByID(ctx context.Context, id string) (*types.Customer, error) {
	var customer *types.Customer

	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no customer found with ID %s", id)
		}
		return nil, fmt.Errorf("error decoding customer: %w", err)
	}

	return customer, nil
}

func (r *Repository) Create(ctx context.Context, customer *types.Customer) (*mongo.InsertOneResult, error) {
	res, err := r.collection.InsertOne(ctx, customer)
	return res, err
}

func (r *Repository) Update(ctx context.Context, id string, customer *types.Customer) error {
	filter := bson.D{{"_id", id}}
	update := bson.M{
		"$set": bson.M{
			"first_name": customer.FirstName,
			"last_name":  customer.LastName,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	filter := bson.D{{"_id", id}}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *Repository) GetCustomersByFilter(ctx context.Context, params model.QueryParams) ([]types.Customer, int64, error) {

	var customers []types.Customer
	limit, offset := pkg.LimitOffsetValidation(params.Limit, params.Offset)

	filter, err := helpers.CreateFilter(params)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid filter parameters"})
	}

	fmt.Printf("Filter: %v\n", filter)

	opts := options.Find().SetLimit(limit).SetSkip(offset)

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error counting customers"})
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Could not get any customers"})
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, 0, echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error decoding customers"})
	}

	return customers, totalCount, nil

}
