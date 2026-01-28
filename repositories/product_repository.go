package repositories

import (
	"context"

	"cashier/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository interface {
	FindAll() ([]model.ProductResponse, error)
	FindByID(id int) (model.ProductResponse, error)
	Save(product model.Product) (model.Product, error)
	Update(product model.Product) (model.Product, error)
	Delete(id int) error
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll() ([]model.ProductResponse, error) {
	var products []model.ProductResponse
	rows, err := r.db.Query(context.Background(), "SELECT p.id, p.name, p.price, p.stock, c.name as category_name FROM products p JOIN categories c ON p.category_id = c.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.ProductResponse
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryName); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) FindByID(id int) (model.ProductResponse, error) {
	var product model.ProductResponse
	err := r.db.QueryRow(context.Background(), "SELECT p.id, p.name, p.price, p.stock, c.name as category_name FROM products p JOIN categories c ON p.category_id = c.id WHERE p.id = $1", id).Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.CategoryName)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) Save(product model.Product) (model.Product, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id", product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) Update(product model.Product) (model.Product, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5", product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return product, err
	}
	return product, nil
}

func (r *productRepository) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM products WHERE id = $1", id)
	return err
}
