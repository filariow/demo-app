package persistence

import (
	"context"
	"eshop-catalog/pkg/models"
)

type Repository interface {
	Read(context.Context, string) (*models.Product, error)
	List(context.Context) ([]models.Product, error)
	AddOrderedUnits(context.Context, string, string, int64) error
	Close(context.Context) error
}
