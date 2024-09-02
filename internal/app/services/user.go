package services

import (
	"context"
	"errors"
	"fmt"
	"motorbike-rental-backend/internal/app/models"
	"motorbike-rental-backend/pkg/errorsx"

	"gorm.io/gorm"
)

type IUserService interface {
	CreateUser(ctx context.Context, user models.User) error
	GetAllUser(ctx context.Context) (*[]models.User, error)
	GetByUserID(ctx context.Context, param int64) (*models.User, error)
	DeleteByUserID(ctx context.Context, param int64) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	MeUpdate(ctx context.Context, m models.User) error
	UpdateUser(ctx context.Context, m models.User) error
}

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) IUserService {
	return &UserService{DB: db}
}

func (u *UserService) CreateUser(ctx context.Context, user models.User) error {
	// Gerekirse context ile bağlantı yönetimi veya zaman aşımı işlemlerini ekleyebiliriz ileride
	return u.DB.WithContext(ctx).Create(&user).Error
}

func (u *UserService) GetAllUser(ctx context.Context) (*[]models.User, error) {
	var users []models.User
	if err := u.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (u *UserService) GetByUserID(ctx context.Context, param int64) (*models.User, error) {
	var user models.User

	if err := u.DB.WithContext(ctx).Where("id = ?", param).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserService) DeleteByUserID(ctx context.Context, param int64) error {
	var user models.User

	if err := u.DB.WithContext(ctx).Where("id = ?", param).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	// Kullanıcı mevcutsa sil
	if err := u.DB.WithContext(ctx).Where("id = ?", param).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	// GORM ile sorgu yapma
	result := u.DB.WithContext(ctx).Where("email = ?", email).Find(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Kayıt bulunamadı hatasını anlamlı bir mesajla döndür
			return nil, gorm.ErrRecordNotFound
		}
		// Diğer veritabanı hatalarını anlamlı bir mesajla döndür
		return nil, fmt.Errorf("veritabanı işlemi sırasında bir hata oluştu: %w", result.Error)
	}

	return &user, nil
}

func (u UserService) MeUpdate(ctx context.Context, m models.User) error {
	var count int64
	err := u.DB.Model(&models.User{}).Where("id != ? AND email = ?", m.ID, m.Email).Count(&count).Error
	if err != nil {
		return errorsx.Database(err)
	}
	if count > 0 {
		return errorsx.BadRequestError("Email başka bir kullanıcı tarafından kullanılmaktadır")
	}

	err = u.DB.WithContext(ctx).Model(&m).Where("id = ?", m.ID).Updates(m).Error
	if err != nil {
		return errorsx.Database(err)
	}

	return nil
}

func (u UserService) UpdateUser(ctx context.Context, m models.User) error {
	var count int64

	// Kendi ID'si dışındaki aynı kullanıcı adı var mı?
	err := u.DB.Model(&models.User{}).
		Where("id != ? AND username = ?", m.ID, m.UserName).
		Count(&count).Error
	if err != nil {
		return errorsx.Database(err)
	}
	if count > 0 {
		return errorsx.BadRequestError("Kullanıcı adı başka bir kullanıcı tarafından kullanılmaktadır")
	}

	// Kendi ID'si dışındaki aynı email adresi var mı?
	err = u.DB.Model(&models.User{}).
		Where("id != ? AND email = ?", m.ID, m.Email).
		Count(&count).Error
	if err != nil {
		return errorsx.Database(err)
	}
	if count > 0 {
		return errorsx.BadRequestError("Email başka bir kullanıcı tarafından kullanılmaktadır")
	}

	// Kullanıcıyı güncelle
	err = u.DB.WithContext(ctx).Model(&m).Where("id = ?", m.ID).Updates(m).Error
	if err != nil {
		return errorsx.Database(err)
	}

	return nil
}
