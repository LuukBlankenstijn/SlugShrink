package db

import (
	"fmt"
	"time"

	"github.com/LuukBlankenstijn/slugshrink/internal/logging"
	"github.com/LuukBlankenstijn/slugshrink/internal/repo/gorm/models"
	"github.com/LuukBlankenstijn/slugshrink/internal/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open() (*gorm.DB, error) {
	gormLogger := &SlogAdapter{
		Logger: logging.Logger(),
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	}
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
		Logger:         gormLogger,
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.RedirectModel{}, &models.DomainModel{}, &models.AuthConfigModel{}); err != nil {
		return nil, err
	}

	return db, nil
}
