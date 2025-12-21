//go:build !prod

package db

import (
	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("./sqlite_db"), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.RedirectModel{}, models.DomainModel{}, &models.AuthConfigModel{}); err != nil {
		return nil, err
	}
	return db, nil
}
