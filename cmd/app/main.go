package main

import (
	"context"
	"log"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/modules/auth"
	"deployer.com/modules/containers"
	"deployer.com/modules/deployments"
	"deployer.com/modules/domains"
	"deployer.com/modules/projects"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"deployer.com/modules/users"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewFiber() *fiber.App {
	return fiber.New()
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})
	{
		group := api.Group("/auth")
		routes := auth.NewAuthController(auth.NewAuthService(users.NewUsersService(db)))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/secrets")
		routes := secrets.NewSecretsController(&group, secrets.NewSecretsService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/users")
		routes := users.NewUsersController(&group, users.NewUsersService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/servers")
		routes := servers.NewServersController(&group, servers.NewServersService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/containers")
		routes := containers.NewContainersController(&group, containers.NewContainersService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/scripts")
		routes := scripts.NewScriptsController(&group, scripts.NewScriptsService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/domains")
		routes := domains.NewDomainsController(&group, domains.NewDomainsService(db))
		routes.RegisterDomainsRoutes(&group)
	}
	{
		group := api.Group("/sub-domains")
		routes := domains.NewSubDomainsController(&group, domains.NewSubDomainsService(db))
		routes.RegisterSubDomainsRoutes(&group)
	}
	{
		group := api.Group("/deployments")
		routes := deployments.NewDeploymentsController(&group, deployments.NewDeploymentsService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/projects")
		routes := projects.NewProjectsController(&group, projects.NewProjectsService(db))
		routes.RegisterRoutes(&group)
	}
}

func main() {
	app := fx.New(
		fx.Provide(
			NewFiber,
			postgres.NewGormDB,
		),
		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, db *gorm.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := db.AutoMigrate(
						&users.User{},
						&secrets.Secret{},
						&servers.Server{},
						&containers.Container{},
						&scripts.Script{},
						&domains.Domain{},
						&domains.SubDomain{},
						&deployments.Deployment{},
						&projects.Project{},
						&projects.ProjectDeployments{},
					); err != nil {
						log.Fatal("AutoMigrate failed:", err)
					}
					RegisterRoutes(app, db)
					go func() {
						if err := app.Listen(":8080"); err != nil {
							log.Fatal(err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return app.Shutdown()
				},
			})
		}),
	)
	app.Run()
}
