package main

import (
	"context"
	"log/slog"

	"github.com/LuukBlankenstijn/gewish/internal/app"
	"github.com/LuukBlankenstijn/gewish/internal/db"
	"github.com/LuukBlankenstijn/gewish/internal/logging"
	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/transport/api/dashboard"
	"github.com/LuukBlankenstijn/gewish/internal/transport/api/public"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logging.Logger().Debug("No .env file found, using system environment variables", slog.Any("error", err))
	}

	ctx := context.Background()
	database, err := db.Open()
	if err != nil {
		logging.Fatal("failed to open database", slog.Any("error", err))
	}

	redirectRepo := gormrepo.NewRedirectsRepo(database)
	domainRepo := gormrepo.NewDomainsRepo(database)
	authConfigRepo := gormrepo.NewAuthConfigRepo(database)

	redirect := app.NewRedirects(redirectRepo)
	domains := app.NewDomains(domainRepo)
	authConfigs := app.NewAuthConfigs(authConfigRepo)

	dashboardApi := dashboard.NewDashboardApi(redirect, domains, authConfigs)
	publicApi := public.NewPublicApi(redirect, domains)
	g, _ := errgroup.WithContext(ctx)

	g.Go(dashboardApi.Run)
	g.Go(publicApi.Run)

	if err := g.Wait(); err != nil {
		logging.Fatal("server run failed", slog.Any("error", err))
	}
}
