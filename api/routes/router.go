package routes

import (
	"motorbike-rental-backend/internal/app/user-and-auth/handlers"
	"motorbike-rental-backend/internal/app/user-and-auth/services"
	"motorbike-rental-backend/pkg/app"
	"motorbike-rental-backend/pkg/router"
	"time"
)

type IdareRouter struct {
}

func NewIdareRouter() *IdareRouter {
	return &IdareRouter{}
}

func (IdareRouter) RegisterRoutes(app *app.App) {
	userService := services.NewUserService(app.DB)
	userHandler := handlers.NewUserHandler(userService)

	authService := services.NewAuthService(app.DB, app.Cfg.Server.JwtSecret, app.Cfg.Server.JwtAccessTokenExpireMinute*time.Minute, app.Cfg.Server.JwtRefreshTokenExpireHour*time.Hour)
	authHandler := handlers.NewAuthHandler(authService, userService)

	api := app.FiberApp.Group("/api")

	router.Post(api, "/user/create", userHandler.CreateUser) // yeni bir kullanıcı hesap oluşturur
	router.Post(api, "/auth/login", authHandler.Login)
	router.Post(api, "/auth/refresh", authHandler.RefreshToken)

	api.Use(router.JWTMiddleware(app))

	router.Get(api, "/user/me", userHandler.Me)
	router.Get(api, "/users", userHandler.GetAllUsers) // tüm kullanıcıları getirir
	router.Put(api, "/user/me", userHandler.MeUpdate)

	router.Post(api, "/auth/logout", authHandler.Logout)

	// Only admins can access them.
	adminRoutes := api.Group("")
	adminRoutes.Use(router.AdminControlMiddleware)

	router.Post(adminRoutes, "/user/create", userHandler.CreateUser)       // admin yeni kullanıcı ekleyebilir
	router.Post(adminRoutes, "/user/createAdmin", userHandler.CreateAdmin) // admin yeni bir admin ekleyebilir
	router.Post(adminRoutes, "/user/:id", userHandler.DeleteByUserID)      // kullanıcı silme
	router.Get(adminRoutes, "/users/:id", userHandler.GetByUserID)         // belli bir kullanıcıyı getirir
	router.Put(adminRoutes, "/user/update/:id", userHandler.UpdateUserByID)

}
