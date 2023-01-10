CREATE TABLE IF NOT EXISTS products (
    id uuid,
    name VARCHAR(100) NOT NULL,
    photoUrl VARCHAR(500) NULL,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS units_ordered (
    order_id uuid,
    product_id uuid,
    unit_sold bigint,
    PRIMARY KEY(order_id, product_id),
    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
            REFERENCES products(id)
);
