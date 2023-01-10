package models

type Product struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	PhotoURL string `json:"photoUrl"`
	UnitSold int64  `json:"unitSold"`
}
