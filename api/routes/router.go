package routes

import (
	motor_Handler "motorbike-rental-backend/internal/app/motorbike/handlers"
	motor_Service "motorbike-rental-backend/internal/app/motorbike/services"
	ride_Handler "motorbike-rental-backend/internal/app/ride/handlers"
	ride_Service "motorbike-rental-backend/internal/app/ride/services"
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

	motorService := motor_Service.NewMotorService(app.DB)
	motorHandler := motor_Handler.NewMotorHandler(motorService)

	rideService := ride_Service.NewRideService(app.DB)
	rideHandler := ride_Handler.NewRideHandler(rideService)

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
	router.Delete(adminRoutes, "/motorbike/:id", motorHandler.DeleteMotor)
	router.Get(adminRoutes, "/motorbikes", motorHandler.GetAllMotors)
	router.Get(adminRoutes, "/motorbikes/:id", motorHandler.GetMotorByID)
	router.Get(adminRoutes, "/available-motorbikes", motorHandler.GetAvailableMotors)
	router.Get(adminRoutes, "/maintenance-motorbikes", motorHandler.GetMaintenanceMotors)
	router.Get(adminRoutes, "/rented-motorbikes", motorHandler.GetRentedMotors)
	router.Get(adminRoutes, "/motorbike-photos/:id", motorHandler.GetPhotosByID)

	// ride operations
	router.Get(adminRoutes, "/rides", rideHandler.GetAllRides)
	router.Get(adminRoutes, "/rides/:id", rideHandler.GetRideByID)
	router.Post(adminRoutes, "/ride", rideHandler.CreateRide)
}
