package seed

import (
	"eshop-orders/pkg/models"
	"time"
)

var Data = []models.Order{
	{
		ID:   "50fb2861-c47e-47b0-9757-501195100edc",
		Date: time.Now(),
		OrderedProducts: []models.OrderedProduct{
			{
				ID:           "419f9b8a-a76e-4f7d-9e5f-36c460808967",
				Name:         "Binding services on kubernetes",
				PhotoURL:     "https://picsum.photos/id/367/1200",
				UnitsOrdered: 367,
			},
		},
	},
	{
		ID:   "8add0704-e0c0-4fa2-a6ef-823cb2867f44",
		Date: time.Now(),
		OrderedProducts: []models.OrderedProduct{
			{
				ID:           "4a43d4d1-fc63-4060-920e-d7461f3496e8",
				Name:         "Brewing Microservices",
				PhotoURL:     "https://picsum.photos/id/425/1200",
				UnitsOrdered: 425,
			},
		},
	},
	{
		ID:   "ebd93cd7-3337-4a78-8c34-01dfdc576c0a",
		Date: time.Now(),
		OrderedProducts: []models.OrderedProduct{
			{
				ID:           "f6b8bf0d-4d5b-4943-88c8-3ec7069f750b",
				Name:         "More on Binding",
				PhotoURL:     "https://picsum.photos/id/436/1200",
				UnitsOrdered: 436,
			},
		},
	},
	{
		ID:   "d2819db6-9b89-4258-baaf-5cda5a3210b1",
		Date: time.Now(),
		OrderedProducts: []models.OrderedProduct{
			{
				ID:           "91b00e7a-c21f-4135-b3a7-c3801ccb67c3",
				Name:         "Hybrid Cloud",
				PhotoURL:     "https://picsum.photos/id/621/1200",
				UnitsOrdered: 621,
			},
		},
	},
	{
		ID:   "eac8552d-e000-4b36-a96e-a67f86de96fe",
		Date: time.Now(),
		OrderedProducts: []models.OrderedProduct{
			{
				ID:           "5a2d2b66-3ef9-48a8-ac0d-db842da4577f",
				Name:         "Complex architectures",
				PhotoURL:     "https://picsum.photos/id/508/1200",
				UnitsOrdered: 508,
			},
		},
	},
}
