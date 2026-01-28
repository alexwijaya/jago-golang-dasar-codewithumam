package repositories

import (
	"context"

	"cashier/model"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CategoryRepository interface {
	FindAll() ([]model.Category, error)
	FindByID(id int) (model.Category, error)
	Save(category model.Category) (model.Category, error)
	Update(category model.Category) (model.Category, error)
	Delete(id int) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	rows, err := r.db.Query(context.Background(), "SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(id int) (model.Category, error) {
	var category model.Category
	err := r.db.QueryRow(context.Background(), "SELECT id, name, description FROM categories WHERE id = $1", id).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Save(category model.Category) (model.Category, error) {
	err := r.db.QueryRow(context.Background(), "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id", category.Name, category.Description).Scan(&category.ID)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Update(category model.Category) (model.Category, error) {
	_, err := r.db.Exec(context.Background(), "UPDATE categories SET name = $1, description = $2 WHERE id = $3", category.Name, category.Description, category.ID)
	if err != nil {
		return category, err
	}
	return category, nil
}

func (r *categoryRepository) Delete(id int) error {
	_, err := r.db.Exec(context.Background(), "DELETE FROM categories WHERE id = $1", id)
	return err
}
