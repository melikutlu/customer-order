package internal

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/CustomerService/internal/types"
	"tesodev-korpes/shared/model"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetByID(ctx context.Context, id string) (*types.CustomerResponseModel, error) {

	response, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	customer := ToCustomerResponse(response)

	if customer == nil {
		return nil, fmt.Errorf("order not found")
	}

	return customer, nil
}

func (s *Service) Create(ctx context.Context, customerRequestModel types.CustomerRequestModel) (string, error) {

	hashedPassword, err := authentication.HashPassword(customerRequestModel.Password)
	if err != nil {
		return "", err
	}

	customID := uuid.New().String()
	now := time.Now().Local()
	customerRequestModel.CreatedAt = now

	customer := &types.Customer{
		FirstName: customerRequestModel.FirstName,
		LastName:  customerRequestModel.LastName,
		Age:       customerRequestModel.Age,
		Email:     customerRequestModel.Email,
		CreatedAt: customerRequestModel.CreatedAt,
		Id:        customID,
		Username:  customerRequestModel.Username,
		Password:  hashedPassword,
	}

	_, err = s.repo.Create(ctx, customer)
	if err != nil {
		return "", err
	}

	return customID, nil
}

func (s *Service) Update(ctx context.Context, id string, customerUpdateModel *types.CustomerUpdateModel) error {

	customer, err := s.repo.FindByID(ctx, id)

	now := time.Now().Local()

	if err != nil {
		return err
	}

	customer.FirstName = customerUpdateModel.FirstName
	customer.LastName = customerUpdateModel.LastName
	customer.Email = customerUpdateModel.Email
	customer.Username = customerUpdateModel.Username
	customer.UpdatedAt = now

	return s.repo.Update(ctx, id, customer)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetCustomers(ctx context.Context, params model.QueryParams) ([]types.CustomerResponseModel, int64, error) {

	response, totalCount, err := s.repo.GetCustomersByFilter(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	customer := ToCustomerRespList(response)
	return customer, totalCount, nil
}

func (s *Service) GetForLogin(ctx context.Context, id string) (*types.CustomerLoginResponseModel, error) {

	response, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	customer := ToLoginResponse(response)
	return customer, nil
}
