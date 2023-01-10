package persistence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"eshop-orders/pkg/models"

	"github.com/google/uuid"
)

func NewInMemoryRepo() Repository {
	pp := []models.Order{
		{
			ID:   uuid.New().String(),
			Date: time.Now(),
			OrderedProducts: []models.OrderedProduct{
				{
					ID:           uuid.New().String(),
					Name:         "Binding services on kubernetes",
					PhotoURL:     "https://picsum.photos/id/367/1200",
					UnitsOrdered: 367,
				},
			},
		},
		{
			ID:   uuid.New().String(),
			Date: time.Now(),
			OrderedProducts: []models.OrderedProduct{
				{
					ID:           uuid.New().String(),
					Name:         "brewing microservices",
					PhotoURL:     "https://picsum.photos/id/425/1200",
					UnitsOrdered: 425,
				},
			},
		},
		{
			ID:   uuid.New().String(),
			Date: time.Now(),
			OrderedProducts: []models.OrderedProduct{
				{
					ID:           uuid.New().String(),
					Name:         "More on Binding",
					PhotoURL:     "https://picsum.photos/id/436/1200",
					UnitsOrdered: 436,
				},
			},
		},
		{
			ID:   uuid.New().String(),
			Date: time.Now(),
			OrderedProducts: []models.OrderedProduct{
				{
					ID:           uuid.New().String(),
					Name:         "Hybrid Cloud",
					PhotoURL:     "https://picsum.photos/id/621/1200",
					UnitsOrdered: 621,
				},
			},
		},
		{
			ID:   uuid.New().String(),
			Date: time.Now(),
			OrderedProducts: []models.OrderedProduct{
				{
					ID:           uuid.New().String(),
					Name:         "Complex architectures",
					PhotoURL:     "https://picsum.photos/id/508/1200",
					UnitsOrdered: 508,
				},
			},
		},
	}

	m := make(map[string]models.Order, len(pp))
	for _, v := range pp {
		m[v.ID] = v
	}

	return &inMemoryRepo{products: m}
}

type inMemoryRepo struct {
	mux      sync.RWMutex
	products map[string]models.Order
}

func (r *inMemoryRepo) Read(ctx context.Context, id string) (*models.Order, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("invalid id %s: %w", id, err)
	}

	r.mux.RLock()
	defer r.mux.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return nil, nil
	}

	return &p, nil
}

func (r *inMemoryRepo) List(context.Context) ([]models.Order, error) {
	r.mux.RLock()
	defer r.mux.RUnlock()

	oo := make([]models.Order, len(r.products))
	counter := 0
	for _, v := range r.products {
		oo[counter] = v
		counter++
	}
	return oo, nil
}

func (r *inMemoryRepo) Create(ctx context.Context, order models.Order) (*models.Order, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	order.ID = uuid.New().String()
	order.Date = time.Now()
	r.products[order.ID] = order

	return &order, nil
}
