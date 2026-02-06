package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (s *ProductService) GetAll(name string) ([]models.Product, error) {
	return s.Repo.GetAll(name)
}

func (s *ProductService) GetByID(id int) (models.Product, error) {
	return s.Repo.GetByID(id)
}

func (s *ProductService) Create(product models.Product) (models.Product, error) {
	return s.Repo.Create(product)
}

func (s *ProductService) Update(product models.Product) (models.Product, error) {
	return s.Repo.Update(product)
}

func (s *ProductService) Delete(id int) error {
	return s.Repo.Delete(id)
}
