package router

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/log"
	"motorbike-rental-backend/pkg/viewmodel"
)

func Get(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Get(path, ctxWrap(h))
}

func Post(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Post(path, ctxWrap(h))
}

func Put(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Put(path, ctxWrap(h))
}

func Delete(r fiber.Router, path string, h func(ctx *app.Ctx) error) {
	r.Delete(path, ctxWrap(h))
}

func ctxWrap(h func(ctx *app.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Get("requestid", "")
		logger := log.GetLogger(rid, zap.String("HTTP_METHOD", c.Method()), zap.String("HTTP_PATH", c.Path()))
		c.Locals("logger", &logger)
		return h(&app.Ctx{Ctx: c})
	}
}

func JWTErrorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Status(fiber.StatusUnauthorized)

	if err.Error() == "Missing or malformed JWT" {
		return ctx.JSON(viewmodel.ResponseModel{ErrorMessage: "token malformed"})
	}
	return ctx.JSON(viewmodel.ResponseModel{ErrorMessage: "token expired"})
}

func JWTMiddleware(app *app.App) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(app.Cfg.Server.JwtSecret),
		ErrorHandler: JWTErrorHandler,
	})
}

func AdminControlMiddleware(c *fiber.Ctx) error {
	// JWT'den kullanıcı bilgisini çekiyoruz
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	// Rol bilgisini çekiyoruz
	role, ok := claims["role"].(float64) // claims'den role'u float64 olarak çekiyoruz

	// Role bilgisi ile ilgili durumları yazdırıyoruz --> these are used for debuging, its okey now!
	// fmt.Println("Token claims:", claims)
	// fmt.Println("Role (float64 olarak):", role)

	if !ok {
		// Eğer role alanı mevcut değilse veya float64 değilse bir hata döndürüyoruz
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Yetkiniz bulunmamaktadır",
		})
	}

	if role != 10 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Yetkiniz yok",
		})
	}

	// Eğer rol 10 ise (admin), bir sonraki handler'a geçilir
	return c.Next()
}
