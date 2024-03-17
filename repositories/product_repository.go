package repositories

import (
	"database/sql"
	"marketplace-app/entities"
	"marketplace-app/utils"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) FindByID(id string) (*entities.Product, error) {
	// language=sql
	query := `
		SELECT id, name, price, image_url, stock, condition, tags, is_purchasable, user_id
		FROM products
		WHERE id = $1
		AND deleted_at IS NULL
	`

	row := r.db.QueryRow(query, id)

	product := &entities.Product{}
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition, &product.Tags, &product.IsPurchasable, &product.UserID)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) Search(filter entities.ProductFilter) ([]entities.Product, utils.MetaResponse, error) {
	query := `
		WITH filter AS (
			SELECT
				$1::boolean AS user_only,
				COALESCE($2, 10)::int AS custom_limit,
				COALESCE($3, 0)::int AS custom_offset,
				$4::text[] AS tags,
				$5::text AS condition,
				$6::boolean AS show_empty_stock,
				$7::numeric AS max_price,
				$8::numeric AS min_price,
				$9::text AS sort_by,
				$10::text AS order_by,
				$11::text AS search,
				NULLIF($12, '')::uuid AS user_id
		)
	`

	queryCount := query + `
		SELECT count(*)
		FROM products
	 `

	query = query + `
		SELECT id, name, price, image_url, stock, products.condition, products.tags, is_purchasable
		FROM products`

	queryJoinAndWhere := `
		CROSS JOIN filter
		WHERE deleted_at IS NULL
		AND CASE WHEN filter.user_only IS NOT NULL THEN 
		    CASE WHEN filter.user_only = true THEN products.user_id = filter.user_id ELSE true END
		ELSE true END
		AND CASE WHEN filter.tags != '{}' THEN products.tags && filter.tags ELSE true END
		AND CASE WHEN filter.condition IS NOT NULL THEN products.condition = filter.condition ELSE true END
		AND CASE WHEN filter.show_empty_stock IS NOT NULL THEN 
		    CASE WHEN filter.show_empty_stock = true THEN true ELSE products.stock > 0 END
		ELSE true END
		AND CASE WHEN filter.max_price IS NOT NULL THEN products.price <= filter.max_price ELSE true END
		AND CASE WHEN filter.min_price IS NOT NULL THEN products.price >= filter.min_price ELSE true END
		AND CASE WHEN filter.search IS NOT NULL THEN products.name ILIKE '%' || filter.search || '%' ELSE true END
	`

	queryCount = queryCount + queryJoinAndWhere
	var total int
	err := r.db.QueryRow(queryCount, filter.UserOnly, filter.Limit, filter.Offset, filter.Tags, filter.Condition, filter.ShowEmptyStock, filter.MaxPrice, filter.MinPrice, filter.SortBy, filter.OrderBy, filter.Search, filter.UserID).Scan(&total)
	if err != nil {
		return nil, utils.MetaResponse{}, err
	}

	var meta utils.MetaResponse
	meta.Total = total
	limit := 10
	if filter.Limit != nil {
		limit = *filter.Limit
	}
	meta.Limit = limit
	offset := 0
	if filter.Offset != nil {
		offset = *filter.Offset
	}
	meta.Offset = offset

	query = query + queryJoinAndWhere
	query = query + `
		ORDER BY CASE WHEN filter.sort_by = 'price' THEN cast(products.price as varchar) ELSE cast(products.created_at at time zone 'UTC' as varchar) END,
		CASE WHEN filter.order_by = 'asc' THEN 1 ELSE -1 END
		LIMIT (SELECT custom_limit FROM filter)
		OFFSET (SELECT custom_offset FROM filter)
	`
	rows, err := r.db.Query(query, filter.UserOnly, filter.Limit, filter.Offset, filter.Tags, filter.Condition, filter.ShowEmptyStock, filter.MaxPrice, filter.MinPrice, filter.SortBy, filter.OrderBy, filter.Search, filter.UserID)
	if err != nil {
		return nil, utils.MetaResponse{}, err
	}
	defer rows.Close()

	products := []entities.Product{}
	for rows.Next() {
		var product entities.Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition, &product.Tags, &product.IsPurchasable)
		if err != nil {
			return nil, utils.MetaResponse{}, err
		}

		products = append(products, product)
	}

	return products, meta, nil
}

func (r *ProductRepository) Create(product *entities.Product) error {
	// language=sql
	query := `
		INSERT INTO products (name, price, image_url, stock, condition, tags, is_purchasable, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query, product.Name, product.Price, product.ImageURL, product.Stock, product.Condition, product.Tags, product.IsPurchasable, product.UserID)
	return err
}

func (r *ProductRepository) Update(product *entities.Product) error {
	// language=sql
	query := `
		UPDATE products
		SET name = $1, price = $2, image_url = $3, condition = $4, tags = $5, is_purchasable = $6, updated_at = now()
		WHERE id = $8
	`

	_, err := r.db.Exec(query, product.Name, product.Price, product.ImageURL, product.Condition, product.Tags, product.IsPurchasable, product.ID)
	return err
}

func (r *ProductRepository) UpdateStock(product *entities.Product) error {
	// language=sql
	query := `
		UPDATE products
		SET stock = $1, updated_at = now()
		WHERE id = $2
	`

	_, err := r.db.Exec(query, product.Stock, product.ID)
	return err
}

func (r *ProductRepository) Delete(id string) error {
	// language=sql
	query := `
		UPDATE products
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)
	return err
}
