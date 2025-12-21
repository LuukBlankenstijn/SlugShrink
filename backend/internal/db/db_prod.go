//go:build prod

package db

import (
	"fmt"
	"net/url"

	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"github.com/LuukBlankenstijn/gewish/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(buildPostgresDSN()), &gorm.Config{
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

func buildPostgresDSN() string {
	user := utils.EnvOrPanic("DB_USER")
	pass := utils.EnvOrPanic("DB_PASSWORD")
	host := utils.EnvOrPanic("DB_HOST")
	port := utils.EnvOrPanic("DB_PORT")
	name := utils.EnvOrPanic("DB_NAME")

	sslMode := utils.EnvOrDefault("DB_SSLMODE", "disable")

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pass),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   name,
	}

	q := u.Query()
	q.Set("sslmode", sslMode)
	u.RawQuery = q.Encode()

	return u.String()
}
