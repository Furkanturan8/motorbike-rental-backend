package handlers

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"motorbike-rental-backend/internal/app/user-and-auth/models"
	"motorbike-rental-backend/internal/app/user-and-auth/services"
	"motorbike-rental-backend/internal/app/user-and-auth/viewmodels"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/errorsx"
	"motorbike-rental-backend/pkg/utils"

	"github.com/google/uuid"
	"strings"
)

type AuthHandler struct {
	authService services.IAuthService
	userService services.IUserService
}

func NewAuthHandler(s services.IAuthService, us services.IUserService) AuthHandler {
	h := AuthHandler{
		authService: s,
		userService: us,
	}

	return h
}

func (h AuthHandler) Login(ctx *app.Ctx) error {
	var vm viewmodel.AuthLoginVM

	// POST isteğinden gelen verileri al
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	user, err := h.userService.GetByEmail(ctx.Context(), utils.EmailTemizle(vm.Email))
	if err != nil {
		return err
	}

	ok := utils.CheckPasswordHash(strings.TrimSpace(vm.Password), user.Password)
	if !ok {
		return errorsx.UnauthorizedError("Hatalı Email veya Parola")
	}

	refreshTokenID := uuid.New()
	tokens, err := h.authService.GenerateTokenPair(user.ID, refreshTokenID, float64(user.Role))
	if err != nil {
		return errorsx.InternalError(err)
	}
	// fmt.Println(user.ID, "bu user id") --> we used for debug
	err = h.authService.CreateAuthRefreshToken(ctx.Context(), refreshTokenID, user.ID, float64(user.Role))
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return ctx.SuccessResponse(result)
}

func (h AuthHandler) LoginAdminPanel(ctx *app.Ctx) error {
	var vm viewmodel.AuthLoginVM

	// POST isteğinden gelen verileri al
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	user, err := h.userService.GetByEmail(ctx.Context(), utils.EmailTemizle(vm.Email))
	if err != nil {
		return err
	}

	if user.Role != models.UserRoleAdmin {
		return errorsx.BadRequestError("Yetkisiz Giriş Denemesi! Yalnızca Adminler Girebilir!")
	}

	ok := utils.CheckPasswordHash(strings.TrimSpace(vm.Password), user.Password)
	if !ok {
		return errorsx.UnauthorizedError("Hatalı Email veya Parola")
	}

	refreshTokenID := uuid.New()
	tokens, err := h.authService.GenerateTokenPair(user.ID, refreshTokenID, float64(user.Role))
	if err != nil {
		return errorsx.InternalError(err)
	}
	// fmt.Println(user.ID, "bu user id") --> we used for debug
	err = h.authService.CreateAuthRefreshToken(ctx.Context(), refreshTokenID, user.ID, float64(user.Role))
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
	return ctx.SuccessResponse(result)
}

func (h AuthHandler) RefreshToken(ctx *app.Ctx) error {
	var vm viewmodel.AuthRefreshVM
	if err := ctx.BodyParseValidate(&vm); err != nil {
		return errorsx.ValidationError(err)
	}

	refreshTokenID, userID, role, err := h.authService.ParseRefreshToken(vm.RefreshToken)
	if err != nil {
		return errorsx.UnauthorizedError(err.Error())
	}
	authRefreshToken, err := h.authService.GetAuthRefreshToken(ctx.Context(), refreshTokenID)
	if err != nil {
		return err
	}
	newTokenPair, err := h.authService.GenerateTokenPair(userID, refreshTokenID, role)
	if err != nil {
		return errorsx.InternalError(err)
	}

	err = h.authService.UpdateAuthRefreshTokenExpires(ctx.Context(), authRefreshToken)
	if err != nil {
		return err
	}

	result := viewmodel.AuthTokenVM{
		AccessToken:  newTokenPair.AccessToken,
		RefreshToken: newTokenPair.RefreshToken,
	}
	return ctx.SuccessResponse(result)
}
func (h AuthHandler) Logout(ctx *app.Ctx) error {
	userID := ctx.GetUserID() // Bu işlevin kullanıcının ID'sini döndürdüğünden emin olun

	if userID == 0 {
		return errors.New("kullanıcı ID'si alınamadı")
	}

	err := h.authService.DeleteAuthRefreshToken(ctx.Context(), userID)
	if err != nil {
		return err
	}

	return nil
}

func (h AuthHandler) CheckRole(ctx *app.Ctx) int64 {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Rol bilgisini çekiyoruz
	role, _ := claims["role"].(float64) // claims'den role'u float64 olarak çekiyoruz

	return int64(role)
}
