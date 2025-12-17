package db

import (
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.RedirectModel{}, models.DomainModel{}); err != nil {
		return nil, err
	}
	return db, nil
}
