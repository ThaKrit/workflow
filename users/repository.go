package users

import "go-api/model"
import "gorm.io/gorm"

type UserRepository interface {
	GetUsers() ([]model.User, error)
	GetUserByID(id uint) (model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user model.User) error
	DeleteUser(user model.User) error
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r *Repository) GetUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) GetUserByID(id uint) (model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}
//ทำ LOGIN
func (r *Repository) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}
//////

func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user model.User) error {
	return r.db.Save(&user).Error
}

func (r *Repository) DeleteUser(user model.User) error {
	return r.db.Delete(&user).Error
}
