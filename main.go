package main

import (
	"Y_scrap/models"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	//ENV
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	//GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Login, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	//Migration
	if err = db.AutoMigrate(&models.Book{}, &models.BookChar{}); err != nil {
		return
	}

	sc := scraper{
		db:  db,
		cfg: cfg,
	}

	if err := sc.collectData(); err != nil {
		fmt.Printf("some err has happened %s\n", err)
	}

	fmt.Println("Completed successfully!")
}
