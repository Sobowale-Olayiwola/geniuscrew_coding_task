package main

import (
	"geniuscrew/domain"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataSources struct {
	MySQLDB *gorm.DB
}

// InitDS establishes connections to fields in dataSources
func initDS() (*DataSources, error) {
	log.Printf("Initializing data sources\n")
	// Initialize MySQLDB connection
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") + "@tcp" + "(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?" + "charset=utf8mb4&parseTime=True&loc=Local"
	mysql.Open(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&domain.Book{}, &domain.Author{})

	return &DataSources{
		MySQLDB: db,
	}, nil
}
