package db

import (
	"fmt"

	"github.com/LuukBlankenstijn/slugshrink/internal/utils"
)

func buildPostgresDSN() string {
	user := utils.EnvOrPanic("DB_USER")
	pass := utils.EnvOrPanic("DB_PASSWORD")
	host := utils.EnvOrPanic("DB_HOST")
	port := utils.EnvOrPanic("DB_PORT")
	name := utils.EnvOrPanic("DB_NAME")
	sslEnabled := utils.EnvOrDefault("DB_SSLMODE", "false") == "true"
	sslMode := "disable"
	if sslEnabled {
		sslMode = "require"
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, pass, name, port, sslMode)
}

func buildMysqlDSN() string {
	user := utils.EnvOrPanic("DB_USER")
	pass := utils.EnvOrPanic("DB_PASSWORD")
	host := utils.EnvOrPanic("DB_HOST")
	port := utils.EnvOrPanic("DB_PORT")
	name := utils.EnvOrPanic("DB_NAME")
	sslEnabled := utils.EnvOrDefault("DB_SSLMODE", "false") == "true"

	baseDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&allowCleartextPasswords=true",
		user, pass, host, port, name)

	if sslEnabled {
		return fmt.Sprintf("%s&tls=true", baseDSN)
	}
	return baseDSN
}

func buildSqliteDSN() string {
	path := utils.EnvOrDefault("DB_PATH", "local.db")
	return fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", path)
}
