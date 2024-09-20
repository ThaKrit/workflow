package items

import "go-api/model"
import "gorm.io/gorm"

type ItemRepository interface {
	GetItems() ([]model.Item, error)
	GetItemByID(id uint) (model.Item, error)
	CreateItem(item *model.Item) error
	UpdateItem(item model.Item) error
	DeleteItem(item model.Item) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r *Repository) GetItems() ([]model.Item, error) {
	var items []model.Item
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) GetItemByID(id uint) (model.Item, error) {
	var item model.Item
	if err := r.db.First(&item, id).Error; err != nil {
		return item, err
	}
	return item, nil
}

func (r *Repository) CreateItem(item *model.Item) error {
	return r.db.Create(item).Error
}

func (r *Repository) UpdateItem(item model.Item) error {
	return r.db.Save(&item).Error
}

func (r *Repository) DeleteItem(item model.Item) error {
	return r.db.Delete(&item).Error
}
