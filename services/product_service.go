package services

import (
	"cashier/model"
	"cashier/repositories"
)

type ProductService interface {
	GetProductsWithFilters(filter model.ProductFilter) ([]model.ProductResponse, error)
	GetProductByID(id int) (model.ProductResponse, error)
	CreateProduct(product model.Product) (model.Product, error)
	UpdateProduct(id int, product model.Product) (model.Product, error)
	DeleteProduct(id int) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetProductsWithFilters(filter model.ProductFilter) ([]model.ProductResponse, error) {
	return s.repo.FindWithFilters(filter)
}

func (s *productService) GetProductByID(id int) (model.ProductResponse, error) {
	return s.repo.FindByID(id)
}

func (s *productService) CreateProduct(product model.Product) (model.Product, error) {
	return s.repo.Save(product)
}

func (s *productService) UpdateProduct(id int, product model.Product) (model.Product, error) {
	product.ID = id
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id int) error {
	return s.repo.Delete(id)
}
