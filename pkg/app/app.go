package app

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"motorbike-rental-backend/pkg/config"

	"motorbike-rental-backend/pkg/database"
	"motorbike-rental-backend/pkg/log"
	"motorbike-rental-backend/pkg/viewmodel"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/fiber/v2/middleware/recover"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type IRouter interface {
	RegisterRoutes(app *App)
}

type App struct {
	FiberApp *fiber.App
	DB       *gorm.DB
	Cfg      *config.Config
	Ctx      context.Context
}

func New(router IRouter, Version, BuildTime string) *App {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	fiberApp := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
		// ErrorHandler: errorsx.ErrorHandler,
		BodyLimit: 20 * 1024 * 1024,
		//ReadBufferSize: fiber.DefaultReadBufferSize * 2, // Request Header Fields Too Large hatası için
	})

	fiberApp.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	fiberApp.Use(cors.New())
	fiberApp.Use(logger.New())
	fiberApp.Use(requestid.New(requestid.Config{
		Header: fiber.HeaderXRequestID,
	}))

	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	fiberApp.Get("/api/version", func(c *fiber.Ctx) error {
		return c.SendString("version: " + Version + " - buildtime: " + BuildTime)
	})

	fiberApp.Use(func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			return err
		}

		if len(c.Response().Body()) == 0 {
			return c.JSON(&viewmodel.ResponseModel{})
		}
		return nil
	})

	db := database.ConnectDB(cfg.Database)

	app := &App{
		FiberApp: fiberApp,
		DB:       db,
		Cfg:      cfg,
		Ctx:      context.Background(),
	}

	router.RegisterRoutes(app)

	return app
}

var l = log.GetLogger("") // loggerımızı tanımladık

func (a *App) MigrateDB() {
	err := database.Migrations(a.DB)
	if err != nil {
		l.Error("DB failed to migrate.")
	}
	fmt.Println("DB migration successfully completed!")
}

func (a *App) Start() {
	l.SetOptions(zap.AddCallerSkip(-2))
	l.Info("http server başlatılıyor...")

	go func() {
		err := a.FiberApp.Listen(fmt.Sprintf(":%v", a.Cfg.Server.Port))
		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	l.Info("Gracefully shutting down...")
	err := a.FiberApp.Shutdown()
	if err != nil {
		l.Error("FiberApp shutdown", zap.Error(err))
	}

	// Veritabanı bağlantısını kapatma
	sqlDB, _ := a.DB.DB()
	if sqlDB != nil {
		err = sqlDB.Close()
		if err != nil {
			l.Error("DB close error", zap.Error(err))
		}
	}

	l.Info("Elahamdülillah")
}
