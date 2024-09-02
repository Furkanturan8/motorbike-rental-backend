package database

import (
	"fmt"
	"motorbike-rental-backend/internal/app/user-and-auth/models"
	"motorbike-rental-backend/pkg/config"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"time"
)

func ConnectDB(cfg config.DbConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbHost, cfg.DbPort, cfg.DbUsername, cfg.DbPassword, cfg.DbName)

	// PostgreSQL veritabanına bağlantı aç
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Gerekirse log seviyesini ayarlayın
	})
	if err != nil {
		log.Errorf("error opening database connection: %v", err)
	}

	// *sql.DB bağlantısını al
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("error getting sql.DB from gorm.DB: %v", err)
	}

	// Veritabanı bağlantı ayarları
	_MaxIdleConn, _ := strconv.Atoi(cfg.MaxIdleConn)
	_MaxPoolSize, _ := strconv.Atoi(cfg.MaxPoolSize)
	_MaxLifetime, _ := strconv.Atoi(cfg.MaxLifetime)

	sqlDB.SetMaxIdleConns(_MaxIdleConn)
	sqlDB.SetMaxOpenConns(_MaxPoolSize)
	sqlDB.SetConnMaxLifetime(time.Duration(_MaxLifetime) * time.Second)

	fmt.Printf("Connected to PostgreSQL database: %s\n", cfg.DbName)

	return db
}

func Migrations(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.AuthTokenPair{},
	)
	if err != nil {
		return fmt.Errorf("AutoMigrate failed: %v", err)
	}
	return nil
}
