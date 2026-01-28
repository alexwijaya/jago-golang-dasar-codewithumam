package services

import (
	"cashier/model"
	"cashier/repositories"
)

type CategoryService interface {
	GetAllCategories() ([]model.Category, error)
	GetCategoryByID(id int) (model.Category, error)
	CreateCategory(category model.Category) (model.Category, error)
	UpdateCategory(id int, category model.Category) (model.Category, error)
	DeleteCategory(id int) error
}

type categoryService struct {
	repo repositories.CategoryRepository
}

func NewCategoryService(repo repositories.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetAllCategories() ([]model.Category, error) {
	return s.repo.FindAll()
}

func (s *categoryService) GetCategoryByID(id int) (model.Category, error) {
	return s.repo.FindByID(id)
}

func (s *categoryService) CreateCategory(category model.Category) (model.Category, error) {
	return s.repo.Save(category)
}

func (s *categoryService) UpdateCategory(id int, category model.Category) (model.Category, error) {
	category.ID = id
	return s.repo.Update(category)
}

func (s *categoryService) DeleteCategory(id int) error {
	return s.repo.Delete(id)
}
