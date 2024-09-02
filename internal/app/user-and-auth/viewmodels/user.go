package viewmodel

import (
	"motorbike-rental-backend/internal/app/user-and-auth/models"
	"motorbike-rental-backend/pkg/utils"
)

type UserLoginVM struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// rolu create ederken set ediyorum gerek yok
type UserCreateVM struct {
	Email    string `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Name     string `json:"name" validate:"required,max=100"`
	Surname  string `json:"surname" validate:"required,max=100"`
	UserName string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,min=3,max=100"`
}

func (vm UserCreateVM) ToDBModel(m models.User) models.User {
	m.Email = utils.EmailTemizle(vm.Email)
	m.Phone = utils.TelefonTemizle(vm.Phone)
	m.Name = utils.ToTitle(vm.Name)
	m.Surname = utils.ToTitle(vm.Surname)
	m.UserName = vm.UserName
	m.Password, _ = utils.HashPassword(vm.Password)

	return m
}

type UserUpdateVM struct {
	Email    string          `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string          `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Name     string          `json:"name" validate:"required,max=100"`
	Surname  string          `json:"surname" validate:"required,max=100"`
	UserName string          `json:"username" validate:"required,max=20"`
	Role     models.UserRole `json:"role" validate:"required"`
	Password string          `json:"password" validate:"max=100"`
}

func (vm UserUpdateVM) ToDBModel(m models.User) models.User {
	m.Email = utils.EmailTemizle(vm.Email)
	m.Phone = utils.TelefonTemizle(vm.Phone)
	m.Name = utils.ToTitle(vm.Name)
	m.Surname = utils.ToTitle(vm.Surname)
	m.UserName = vm.UserName
	m.Role = vm.Role
	if vm.Password != "" {
		m.Password, _ = utils.HashPassword(vm.Password)
	}

	return m
}

type UserListVM struct {
	ID       int64           `json:"id"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	UserName string          `json:"username"`
	Role     models.UserRole `json:"role"`
}

func (vm UserListVM) ToViewModel(m models.User) UserListVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.UserName = m.UserName
	vm.Role = m.Role

	return vm
}

type UserDetailVM struct {
	ID       int64           `json:"id"`
	Email    string          `json:"email"`
	Phone    string          `json:"phone"`
	Name     string          `json:"name"`
	Surname  string          `json:"surname"`
	UserName string          `json:"username"`
	Role     models.UserRole `json:"role"`
}

func (vm UserDetailVM) ToViewModel(m models.User) UserDetailVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.UserName = m.UserName
	vm.Role = m.Role

	return vm
}

type UserMeVM struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	UserName string `json:"username"`
	Role     string `json:"role"`
}

func (vm UserMeVM) ToViewModel(m models.User) UserMeVM {
	vm.ID = m.ID
	vm.Email = m.Email
	vm.Phone = m.Phone
	vm.Name = m.Name
	vm.Surname = m.Surname
	vm.UserName = m.UserName
	vm.Role = m.Role.String()

	return vm
}

type UserMeUpdateVM struct {
	Email    string `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Name     string `json:"name" validate:"required,max=100"`
	Surname  string `json:"surname" validate:"required,max=100"`
	UserName string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"max=100"`
}

func (vm UserMeUpdateVM) ToDBModel(m models.User) models.User {
	m.Email = utils.EmailTemizle(vm.Email)
	m.Phone = utils.TelefonTemizle(vm.Phone)
	m.Name = utils.ToTitle(vm.Name)
	m.Surname = utils.ToTitle(vm.Surname)
	m.UserName = vm.UserName
	if vm.Password != "" {
		m.Password, _ = utils.HashPassword(vm.Password)
	}

	return m
}
