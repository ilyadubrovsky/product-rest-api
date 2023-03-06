package product

import "context"

type Repository interface {
	Create(ctx context.Context, prdct Product) (string, error)
	FindAll(ctx context.Context) ([]Product, error)
	FindOne(ctx context.Context, id string) (Product, error)
	FullyUpdate(ctx context.Context, id string, prdct Product) error
	PartiallyUpdate(ctx context.Context, id string, prdct Product) error
	Delete(ctx context.Context, id string) error
}
