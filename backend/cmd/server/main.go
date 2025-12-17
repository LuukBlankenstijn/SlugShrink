package main

import (
	"log"

	"github.com/LuukBlankenstijn/gewish/internal/app"
	"github.com/LuukBlankenstijn/gewish/internal/db"
	gormrepo "github.com/LuukBlankenstijn/gewish/internal/repo/gorm"
	"github.com/LuukBlankenstijn/gewish/internal/transport/api"
)

func main() {
	database, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	redirectRepo := gormrepo.NewRedirectsRepo(database)
	domainRepo := gormrepo.NewDomainsRepo(database)

	redirect := app.NewRedirects(redirectRepo)
	domains := app.NewDomains(domainRepo)

	server := api.NewApiServer(redirect, domains)

	log.Fatal(server.Run())
}
