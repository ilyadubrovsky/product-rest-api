package product

import "context"

type Repository interface {
	Create(ctx context.Context, dto CreateProductDTO) (string, error)
	FindAll(ctx context.Context) ([]Product, error)
	FindOne(ctx context.Context, id string) (Product, error)
	FullyUpdate(ctx context.Context, id string, dto FullyUpdateProductDTO) error
	PartiallyUpdate(ctx context.Context, id string, dto PartiallyUpdateProductDTO) error
	Delete(ctx context.Context, id string) error
}
