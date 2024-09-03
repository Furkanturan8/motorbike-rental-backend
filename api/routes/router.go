package routes

import (
	bikeHandler "motorbike-rental-backend/internal/app/motorbike/handlers"
	bikeService "motorbike-rental-backend/internal/app/motorbike/services"
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

	motorService := bikeService.NewMotorService(app.DB)
	motorHandler := bikeHandler.NewMotorHandler(motorService)

	api := app.FiberApp.Group("/api")

	router.Post(api, "/user/create", userHandler.CreateUser)
	router.Post(api, "/auth/login", authHandler.Login)
	router.Post(api, "/auth/refresh", authHandler.RefreshToken)

	api.Use(router.JWTMiddleware(app))

	router.Get(api, "/user/me", userHandler.Me)
	router.Put(api, "/user/me", userHandler.MeUpdate)

	router.Post(api, "/auth/logout", authHandler.Logout)

	// Only admins can access them.
	adminRoutes := api.Group("")
	adminRoutes.Use(router.AdminControlMiddleware)

	// user operations
	router.Get(adminRoutes, "/users", userHandler.GetAllUsers)
	router.Get(adminRoutes, "/users/:id", userHandler.GetByUserID)
	router.Post(adminRoutes, "/user/create", userHandler.CreateUser)
	router.Post(adminRoutes, "/user/createAdmin", userHandler.CreateAdmin)
	router.Post(adminRoutes, "/user/:id", userHandler.DeleteByUserID)
	router.Put(adminRoutes, "/user/update/:id", userHandler.UpdateUserByID)

	// motorbike operations
	router.Post(adminRoutes, "/motorbike", motorHandler.CreateMotor)
	router.Put(adminRoutes, "/motorbike/:id", motorHandler.UpdateMotor)
	router.Get(adminRoutes, "/motorbikes", motorHandler.GetAllMotors)
	router.Get(adminRoutes, "/motorbikes/:id", motorHandler.GetMotorByID)
	router.Get(adminRoutes, "/available-motorbikes", motorHandler.GetAvailableMotors)
	router.Get(adminRoutes, "/maintenance-motorbikes", motorHandler.GetMaintenanceMotors)
	router.Get(adminRoutes, "/rented-motorbikes", motorHandler.GetRentedMotors)
	router.Get(adminRoutes, "/motorbike-photos/:id", motorHandler.GetPhotosByID)
}
