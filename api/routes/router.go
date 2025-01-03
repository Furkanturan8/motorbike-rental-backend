package routes

import (
	_connHandler "motorbike-rental-backend/internal/app/bluetooth-connection/handlers"
	_connService "motorbike-rental-backend/internal/app/bluetooth-connection/services"
	_mapHandler "motorbike-rental-backend/internal/app/map/handlers"
	_mapService "motorbike-rental-backend/internal/app/map/services"
	_motorHandler "motorbike-rental-backend/internal/app/motorbike/handlers"
	_motorService "motorbike-rental-backend/internal/app/motorbike/services"
	_rideHandler "motorbike-rental-backend/internal/app/ride/handlers"
	_rideService "motorbike-rental-backend/internal/app/ride/services"
	_baseHandler "motorbike-rental-backend/internal/app/user-and-auth/handlers"
	_baseService "motorbike-rental-backend/internal/app/user-and-auth/services"
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
	userService := _baseService.NewUserService(app.DB)
	userHandler := _baseHandler.NewUserHandler(userService)

	authService := _baseService.NewAuthService(app.DB, app.Cfg.Server.JwtSecret, app.Cfg.Server.JwtAccessTokenExpireMinute*time.Minute, app.Cfg.Server.JwtRefreshTokenExpireHour*time.Hour)
	authHandler := _baseHandler.NewAuthHandler(authService, userService)

	motorService := _motorService.NewMotorService(app.DB)
	motorHandler := _motorHandler.NewMotorHandler(motorService)

	mapService := _mapService.NewMapService(app.DB)
	mapHandler := _mapHandler.NewMapHandler(mapService, motorService)

	connService := _connService.NewConnService(app.DB)
	connHandler := _connHandler.NewConnHandler(connService, motorService)

	rideService := _rideService.NewRideService(app.DB)
	rideHandler := _rideHandler.NewRideHandler(rideService, motorService, connHandler)

	api := app.FiberApp.Group("/api")

	router.Post(api, "/user/create", userHandler.CreateUser)
	router.Post(api, "/auth/login", authHandler.Login)
	router.Post(api, "/auth/refresh", authHandler.RefreshToken)

	// admin panel login
	router.Post(api, "/auth/admin/login", authHandler.LoginAdminPanel)

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
	router.Delete(adminRoutes, "/user/:id", userHandler.DeleteByUserID)
	router.Put(adminRoutes, "/user/update/:id", userHandler.UpdateUserByID)

	// motorbike operations
	router.Post(adminRoutes, "/motorbike", motorHandler.CreateMotor)
	router.Put(adminRoutes, "/motorbike/:id", motorHandler.UpdateMotor)
	router.Delete(adminRoutes, "/motorbike/:id", motorHandler.DeleteMotor)
	router.Get(api, "/motorbikes", motorHandler.GetAllMotors)
	router.Get(api, "/motorbikes/:id", motorHandler.GetMotorByID)
	router.Get(api, "/available-motorbikes", motorHandler.GetAvailableMotors)
	router.Get(adminRoutes, "/maintenance-motorbikes", motorHandler.GetMaintenanceMotors)
	router.Get(adminRoutes, "/rented-motorbikes", motorHandler.GetRentedMotors)
	router.Get(adminRoutes, "/motorbike-photos/:id", motorHandler.GetPhotosByID)

	// ride operations
	router.Get(adminRoutes, "/rides", rideHandler.GetAllRides)
	router.Get(adminRoutes, "/rides/:id", rideHandler.GetRideByID)
	router.Post(api, "/ride", rideHandler.CreateRide)
	router.Get(adminRoutes, "/rides/user/:userID", rideHandler.GetRidesByUserID) // Belirli bir kullanıcıya ait tüm kiralamaları getirme
	router.Get(api, "/users/:userID/rides/:rideID", rideHandler.GetRideByUserID) // frontend'de getRideByMe olarak sadece her kullanıcının kendi id'leri gitmeli.
	router.Get(adminRoutes, "/motorbike/:bikeID/rides", rideHandler.GetRidesByBikeID)
	router.Put(adminRoutes, "/ride/update/:id", rideHandler.UpdateRideByID)
	router.Delete(adminRoutes, "/ride/:id", rideHandler.DeleteRide)
	router.Get(adminRoutes, "/filtered-rides", rideHandler.GetRidesByDateRange)      // belirli tarih aralıklarındaki sürüşleri getirir -> /filtered-rides?start_time=2024-09-04&end_time=2024-09-05
	router.Get(api, "/rides/user/:userID/filter", rideHandler.GetRidesByUserAndDate) // userID ye göre belirli tarihler arasında getirir -> /rides/user/:userID/filter?start_time=2024-09-01&end_time=2024-09-09
	router.Put(api, "/ride/finish/:id", rideHandler.FinishRide)
	router.Post(api, "/ride/:id/photo", rideHandler.AddRidePhoto)

	// map operations
	router.Post(adminRoutes, "/map", mapHandler.CreateMap)
	router.Delete(adminRoutes, "/map/:id", mapHandler.DeleteMap)
	router.Get(api, "/maps", mapHandler.GetAllMaps)
	router.Get(api, "/maps/:id", mapHandler.GetMapByID)
	router.Get(api, "/motorbikes/:motorbikeID/map", mapHandler.GetMapByMotorID)
	router.Put(adminRoutes, "/map/update/:id", mapHandler.UpdateMap)
	router.Put(adminRoutes, "/motorbikes/:motorbikeID/map/update", mapHandler.UpdateMapByMotorID)

	// bluetooth connection operations
	router.Get(adminRoutes, "/connections", connHandler.GetAllConnections)
	router.Get(api, "/connections/:id", connHandler.GetConnByID)
	router.Get(adminRoutes, "/connection/motorbike/:motorbikeID", connHandler.GetConnByMotorID)
	router.Get(api, "/connection/user/:userID", connHandler.GetConnByUserID)
	router.Post(api, "/connection/connect", connHandler.Connect) // connect
	router.Delete(adminRoutes, "/connection/:id", connHandler.DeleteConn)
	// router.Post(adminRoutes, "/connection/disconnect/:id", connHandler.Disconnect) // disconnect
}

// Sürüşü bitirme işlem süreci:
// önce kullanıcı motoru kitleyecek, (burada özel bir işlem yapılmalı, kitlendiğinde bana bilgi gelmesi gerekiyor! Burada bluetooth cihazı bana info atacak!)
// daha sonra kullanıcı fotoğrafı yükleyecek, disconnect fonksiyonu çalışacak (motor available olacak)
// daha sonra finishRide fonksiyonuna istek atılacak! Ve burada işlem ücreti çıkacak! Buradan payment fonksiyonuna istek atılıp para alınacak!
