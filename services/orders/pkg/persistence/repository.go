package persistence

import (
	"context"

	"eshop-orders/pkg/models"
)

type Repository interface {
	Create(context.Context, models.Order) (*models.Order, error)
	Read(context.Context, string) (*models.Order, error)
	List(context.Context) ([]models.Order, error)
}
