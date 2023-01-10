package models

import "time"

type Order struct {
	ID              string           `json:"id"`
	OrderedProducts []OrderedProduct `json:"orderedProducts"`
	Date            time.Time        `json:"date"`
}

type OrderedProduct struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PhotoURL     string `json:"photoURL"`
	UnitsOrdered int64  `json:"unitsOrdered"`
}
