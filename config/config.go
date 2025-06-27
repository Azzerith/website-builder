package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"website-builder/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	// Load DB config from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	charset := os.Getenv("DB_CHARSET")
	parseTime := os.Getenv("DB_PARSE_TIME")
	loc := os.Getenv("DB_LOC")

	// Format DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		user,
		password,
		host,
		port,
		dbname,
		charset,
		parseTime,
		loc,
	)

	// Konfigurasi logger GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Threshold untuk query slow
			LogLevel:                  logger.Silent, // Log level (Silent, Error, Warn, Info)
			IgnoreRecordNotFoundError: true,          // Ignore error record not found
			Colorful:                  false,         // Disable color
		},
	)

	// Membuka koneksi database dengan konfigurasi
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	DB = db
	log.Println("Database connected successfully!")

	// Auto migrate semua model
	err = db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.TeamMember{},
		&models.Project{},
		&models.Page{},
		&models.Element{},
		&models.Revision{},
		&models.Comment{},
		&models.CommentReply{},
		&models.Session{},
		&models.Template{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate database: ", err)
	}

	log.Println("Database migration completed!")
}

func GetDB() *gorm.DB {
	return DB
}