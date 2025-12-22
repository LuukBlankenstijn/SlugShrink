package db

import (
	"fmt"

	"github.com/LuukBlankenstijn/gewish/internal/repo/gorm/models"
	"github.com/LuukBlankenstijn/gewish/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Open() (*gorm.DB, error) {
	dbType := utils.EnvOrDefault("DB_TYPE", "sqlite") // options: postgres, mysql, sqlite

	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dialector = postgres.Open(buildPostgresDSN())
	case "mysql", "mariadb":
		dialector = mysql.Open(buildMysqlDSN())
	case "sqlite":
		dialector = sqlite.Open(buildSqliteDSN())
	default:
		return nil, fmt.Errorf("unsupported DB_TYPE: %s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.RedirectModel{}, &models.DomainModel{}, &models.AuthConfigModel{}); err != nil {
		return nil, err
	}

	return db, nil
}
