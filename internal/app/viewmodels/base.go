package viewmodel

import (
	"motorbike-rental-backend/internal/app/models"
)

type UserDetailVMForAdmin struct {
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	UserName string          `json:"username"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Role     models.UserRole `json:"role"`
}

type UserDetailVMForUser struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	UserName string `json:"username"`
	Phone    string `json:"phone"`
}

// ha bu team modelinde captain olarak user kullandık fakat listelerken sadece gerekli şeyleri listelesin diye bu fonksiyonu yazdık.
func UserToUserDetailVMForAdmin(user models.User) UserDetailVMForAdmin {
	return UserDetailVMForAdmin{
		Name:     user.Name,
		Surname:  user.Surname,
		UserName: user.UserName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
	}
}

func UserToUserDetailVMForUser(user models.User) UserDetailVMForUser {
	return UserDetailVMForUser{
		Name:     user.Name,
		Surname:  user.Surname,
		UserName: user.UserName,
		Phone:    user.Phone,
	}
}
