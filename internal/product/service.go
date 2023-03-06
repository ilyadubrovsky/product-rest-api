package product

import (
	"context"
	"github.com/ilyadubrovsky/product-rest-api/internal/config"
	"github.com/ilyadubrovsky/product-rest-api/pkg/logging"
	"github.com/pkg/errors"
)

var (
	ErrNotFound       = errors.New("resource not found")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
)

type Service struct {
	logger     *logging.Logger
	cfg        *config.Config
	repository Repository
}

func NewService(logger *logging.Logger, cfg *config.Config, repository Repository) *Service {
	return &Service{logger: logger, cfg: cfg, repository: repository}
}

func (s *Service) GetAllProducts(ctx context.Context) ([]Product, error) {
	prdcts, err := s.repository.FindAll(ctx)
	if err != nil {
		s.logger.Warningf("failed to FindAll due to error: %v", err)
		return nil, err
	}

	return prdcts, nil
}

func (s *Service) GetProductByID(ctx context.Context, id string) (Product, error) {
	prdct, err := s.repository.FindOne(ctx, id)
	if err != nil {
		s.logger.Warningf("failed to FindOne due to error: %v", err)
		return Product{}, err
	}

	return prdct, nil
}

func (s *Service) CreateProduct(ctx context.Context, dto CreateProductDTO) (string, error) {
	createdProduct := NewProduct(dto)

	id, err := s.repository.Create(ctx, createdProduct)
	if err != nil {
		s.logger.Warningf("failed to Create due to error: %v", err)
		return "", err
	}

	return id, nil
}

func (s *Service) FullyUpdateProductByID(ctx context.Context, id string, dto UpdateProductDTO) error {
	updatedProduct := UpdateProduct(dto)

	if err := s.repository.FullyUpdate(ctx, id, updatedProduct); err != nil {
		s.logger.Warningf("failed to FullyUpdate due to error: %v", err)
		return err
	}

	return nil
}

func (s *Service) PartiallyUpdateProductByID(ctx context.Context, id string, dto UpdateProductDTO) error {
	updatedProduct := UpdateProduct(dto)

	if err := s.repository.PartiallyUpdate(ctx, id, updatedProduct); err != nil {
		s.logger.Warningf("failed to PartiallyUpdate due to error: %v", err)
		return err
	}

	return nil
}

func (s *Service) DeleteProductByID(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		s.logger.Warningf("failed to Delete due to error: %v", err)
		return err
	}

	return nil
}
