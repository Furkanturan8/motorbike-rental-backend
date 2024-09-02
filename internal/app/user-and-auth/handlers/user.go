package handlers

import (
	"errors"
	"fmt"
	"motorbike-rental-backend/internal/app/user-and-auth/models"
	"motorbike-rental-backend/internal/app/user-and-auth/services"
	"motorbike-rental-backend/internal/app/user-and-auth/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type UserHandler struct {
	userService services.IUserService
}

func NewUserHandler(s services.IUserService) UserHandler {
	return UserHandler{userService: s}
}

func (h UserHandler) BaseCreateUser(ctx *app.Ctx, role int64) error {
	var vm viewmodel.UserCreateVM
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	user := vm.ToDBModel(models.User{})
	user.Role = models.UserRole(role)

	err := h.userService.CreateUser(ctx.Context(), user)
	if err != nil {
		return err
	}

	return nil
}

func (h UserHandler) CreateUser(ctx *app.Ctx) error {
	err := h.BaseCreateUser(ctx, 1)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı oluşturulurken bir hata oluştu."})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User başarıyla eklendi!"})
}

func (h UserHandler) CreateAdmin(ctx *app.Ctx) error {
	err := h.BaseCreateUser(ctx, 10)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Admin oluşturulurken bir hata oluştu."})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Admin başarıyla eklendi!"})
}

func (h UserHandler) GetAllUsers(ctx *app.Ctx) error {
	// Veritabanından tüm kullanıcıları çekiyoruz
	users, err := h.userService.GetAllUser(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcılar getirilirken bir hata oluştu."})
	}

	// Kullanıcıları view model'e dönüştürüyoruz
	userListVM := make([]viewmodel.UserListVM, len(*users))
	for i, user := range *users {
		userListVM[i] = viewmodel.UserListVM{}.ToViewModel(user)
	}

	// Kullanıcıları JSON formatında geri döndürüyoruz
	return ctx.SuccessResponse(userListVM, len(userListVM))
}

func (h UserHandler) GetByUserID(ctx *app.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz user id!"})
	}

	user, err := h.userService.GetByUserID(ctx.Context(), int64(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Kullanıcı bulunamadı!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı getirilirken hata oluştu!"})
	}

	vm := viewmodel.UserDetailVM{}.ToViewModel(*user)
	return ctx.SuccessResponse(vm, 1)
}

func (h UserHandler) DeleteByUserID(ctx *app.Ctx) error {
	param := ctx.Params("id")

	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz user id!"})
	}
	err = h.userService.DeleteByUserID(ctx.Context(), int64(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Böyle bir kullanıcı zaten yok!"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı silinirken hata oluştu!"})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User başarıyla silindi!"})
}

func (h UserHandler) Me(ctx *app.Ctx) error {
	id := ctx.GetUserID()
	user, err := h.userService.GetByUserID(ctx.Context(), id)
	if err != nil {
		return err
	}

	vm := viewmodel.UserMeVM{}.ToViewModel(*user)

	return ctx.SuccessResponse(vm, 1)
}

func (h UserHandler) MeUpdate(ctx *app.Ctx) error {
	id := ctx.GetUserID()
	fmt.Println("güncellenen id :", id)
	m, err := h.userService.GetByUserID(ctx.Context(), id)
	if err != nil {
		return err
	}

	var vm viewmodel.UserMeUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	updatedUser := vm.ToDBModel(*m)
	err = h.userService.MeUpdate(ctx.Context(), updatedUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User başarıyla güncellendi!"})
}

func (h UserHandler) UpdateUserByID(ctx *app.Ctx) error {
	param := ctx.Params("id")
	id, err := strconv.Atoi(param)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz user id!"})
	}

	m, err := h.userService.GetByUserID(ctx.Context(), int64(id))
	if err != nil {
		return err
	}
	var vm viewmodel.UserUpdateVM
	if errs := ctx.BodyParseValidate(&vm); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	updatedUser := vm.ToDBModel(*m)
	err = h.userService.UpdateUser(ctx.Context(), updatedUser)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User başarıyla güncellendi!"})
}
