package main

import (
	"context"
	"log"

	"github.com/LuukBlankenstijn/gewish/internal/app"
	"github.com/LuukBlankenstijn/gewish/internal/db"
	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/transport/api/dashboard"
	"github.com/LuukBlankenstijn/gewish/internal/transport/api/public"
	"golang.org/x/sync/errgroup"
)

func main() {
	ctx := context.Background()
	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
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

	log.Fatal(g.Wait())
}
