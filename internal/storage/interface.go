package storage

import (
	"context"

	"github.com/raghavendrah25/golang-backend/internal/product"
)

type Storage interface {
	CreateProduct(product product.Product, ctx context.Context) error
	GetProducts(ctx context.Context) ([]product.Product, error)
}
