package service

import (
	"demo-service/internal/model"
	"demo-service/internal/repository"
	"fmt"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) Create(req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Получаем созданный продукт с временными метками
	createdProduct, err := s.productRepo.GetByID(product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created product: %w", err)
	}

	return createdProduct, nil
}

func (s *ProductService) GetByID(id int64) (*model.Product, error) {
	product, err := s.productRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return product, nil
}

func (s *ProductService) List(page, limit int) (*model.ProductListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	products, total, err := s.productRepo.List(page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return &model.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func (s *ProductService) Update(id int64, req *model.UpdateProductRequest) (*model.Product, error) {
	updates := make(map[string]interface{})

	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}

	if len(updates) == 0 {
		return s.GetByID(id)
	}

	if err := s.productRepo.Update(id, updates); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return s.GetByID(id)
}

func (s *ProductService) Delete(id int64) error {
	if err := s.productRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	return nil
}




