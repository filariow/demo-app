package persistence

import (
	"context"
	"eshop-catalog/pkg/models"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type postgresRepo struct {
	conn *pgx.Conn
}

const (
	getProductSQL = `
SELECT p.id,p.name,p.photoUrl,sum(uo.unit_sold)
	FROM products p JOIN units_ordered uo ON p.id = uo.product_id
	WHERE id = $1
	GROUP BY p.id,p.name,p.photoUrl;
`
	listProductsSQL = `
SELECT p.id,p.name,p.photoUrl,sum(uo.unit_sold)
	FROM products p JOIN units_ordered uo ON p.id = uo.product_id
	GROUP BY p.id,p.name,p.photoUrl;
`
	addOrderedUnitsSQL = `INSERT INTO units_ordered VALUES ($1,$2,$3)`
)

func NewPostgresRepo(ctx context.Context, url string) (Repository, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &postgresRepo{conn: conn}, nil
}

func (r *postgresRepo) Read(ctx context.Context, id string) (*models.Product, error) {
	row := r.conn.QueryRow(ctx, getProductSQL, id)

	var p models.Product
	if err := row.Scan(&p.ID, &p.Name, &p.PhotoURL, &p.UnitSold); err != nil {
		return nil, fmt.Errorf("error fetching product with id '%s': %w", id, err)
	}
	return &p, nil
}

func (r *postgresRepo) List(ctx context.Context) ([]models.Product, error) {
	d, err := r.conn.Query(ctx, listProductsSQL)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}

	pp := []models.Product{}
	for d.Next() {
		var p models.Product
		if err := d.Scan(&p.ID, &p.Name, &p.PhotoURL, &p.UnitSold); err != nil {
			return nil, fmt.Errorf("error mapping product data %w", err)
		}
		pp = append(pp, p)
	}

	return pp, nil
}

func (r *postgresRepo) AddOrderedUnits(ctx context.Context, orderId string, productId string, units int64) error {
	if _, err := r.conn.Exec(ctx, addOrderedUnitsSQL, orderId, productId, units); err != nil {
		return fmt.Errorf(
			"error adding ordered units for product '%s' and order '%s': %w",
			productId, orderId, err)
	}
	return nil
}

func (r *postgresRepo) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}
