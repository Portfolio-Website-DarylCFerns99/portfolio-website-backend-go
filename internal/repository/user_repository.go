package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"portfolio-website-backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmailorUsername(email string) (*models.User, error)
	Update(id uuid.UUID, data map[string]interface{}) (*models.User, error)
	List() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// func (r *userRepository) GetByUsername(username string) (*models.User, error) {
// 	var user models.User
// 	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, err
// 	}
// 	return &user, nil
// }

//	func (r *userRepository) GetByEmail(email string) (*models.User, error) {
//		var user models.User
//		if err := r.db.First(&user, "email = ?", email).Error; err != nil {
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				return nil, nil
//			}
//			return nil, err
//		}
//		return &user, nil
//	}
func (r *userRepository) GetByEmailorUsername(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ? or username = ?", email, email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(id uuid.UUID, data map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	if err := r.db.Model(&user).Updates(data).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) List() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
