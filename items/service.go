package items

import "go-api/model"

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s *Service) GetItems() ([]model.Item, error) {
	return s.repo.GetItems()
}

func (s *Service) GetItemByID(id uint) (model.Item, error) {
	return s.repo.GetItemByID(id)
}

func (s *Service) CreateItem(item *model.Item) error {
	return s.repo.CreateItem(item)
}

func (s *Service) UpdateItem(item model.Item) error {
	return s.repo.UpdateItem(item)
}

func (s *Service) DeleteItem(item model.Item) error {
	return s.repo.DeleteItem(item)
}
