package persistence

import (
	"context"
	"eshop-catalog/pkg/models"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func NewInMemoryRepo() Repository {
	pp := []models.Product{
		{
			ID:       "419f9b8a-a76e-4f7d-9e5f-36c460808967",
			Name:     "Binding services on kubernetes",
			PhotoURL: "https://picsum.photos/id/367/1200",
			UnitSold: 367,
		},
		{
			ID:       "4a43d4d1-fc63-4060-920e-d7461f3496e8",
			Name:     "Brewing Microservices",
			PhotoURL: "https://picsum.photos/id/425/1200",
			UnitSold: 425,
		},
		{
			ID:       "f6b8bf0d-4d5b-4943-88c8-3ec7069f750b",
			Name:     "More on Binding",
			PhotoURL: "https://picsum.photos/id/436/1200",
			UnitSold: 425,
		},
		{
			ID:       "91b00e7a-c21f-4135-b3a7-c3801ccb67c3",
			Name:     "Hybrid Cloud",
			PhotoURL: "https://picsum.photos/id/621/1200",
			UnitSold: 425,
		},
		{
			ID:       "5a2d2b66-3ef9-48a8-ac0d-db842da4577f",
			Name:     "Complex architectures",
			PhotoURL: "https://picsum.photos/id/508/1200",
			UnitSold: 425,
		},
	}

	mp := make(map[string]models.Product, len(pp))
	for _, p := range pp {
		mp[p.ID] = p
	}
	return &inMemoryRepo{products: mp}
}

type inMemoryRepo struct {
	mux      sync.RWMutex
	products map[string]models.Product
}

func (r *inMemoryRepo) Read(ctx context.Context, id string) (*models.Product, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("Invalid id %s: %w", id, err)
	}

	r.mux.RLock()
	defer r.mux.RUnlock()

	p := r.products[id]
	return &p, nil
}

func (r *inMemoryRepo) List(context.Context) ([]models.Product, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	pp := make([]models.Product, len(r.products))
	counter := 0
	for _, v := range r.products {
		pp[counter] = v
		counter++
	}
	return pp, nil
}

func (r *inMemoryRepo) AddOrderedUnits(ctx context.Context, productId string, orderId string, orderedUnits int64) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	p, ok := r.products[productId]
	if !ok {
		return fmt.Errorf("product with id '%s' not found", productId)
	}

	p.UnitSold += orderedUnits
	r.products[productId] = p

	return nil
}

func (r *inMemoryRepo) Close(context.Context) error { return nil }
