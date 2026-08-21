package category

import (
	"errors"
)

var (
	ErrCategoryNameEmpty = errors.New("category name is empty")
)

type CategoryRequest struct {
	Category string `json:"category"`
	IsActive bool   `json:"is_active"`
}

type CategoryService struct {
	categoryRepo *Repository
}

func NewService(categoryRepo *Repository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) validateCategoryRequest(req *CategoryRequest) error {
	if req.Category == "" {
		return ErrCategoryNameEmpty
	}
	return nil
}

func (s *CategoryService) CreateCategory(req *CategoryRequest) (*Model, error) {
	if err := s.validateCategoryRequest(req); err != nil {
		return nil, err
	}
	category := &Model{
		Category: req.Category,
		IsActive: req.IsActive,
	}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) GetAllCategories() ([]*Model, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoryService) UpdateCategory(id int64, req *CategoryRequest) (*Model, error) {
	if err := s.validateCategoryRequest(req); err != nil {
		return nil, err
	}

	category := &Model{
		ID:       id,
		Category: req.Category,
		IsActive: req.IsActive,
	}

	if err := s.categoryRepo.Update(id, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(id int64) error {
	err := s.categoryRepo.Delete(id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	return nil
}
