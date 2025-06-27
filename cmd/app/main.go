package main

import (
	"context"
	"log"
	"time"

	postgres "deployer.com/cmd/db/db"
	"deployer.com/libs"
	"deployer.com/modules/auth"
	"deployer.com/modules/containers"
	"deployer.com/modules/deployments"
	"deployer.com/modules/domains"
	"deployer.com/modules/execute"
	"deployer.com/modules/projects"
	"deployer.com/modules/scripts"
	"deployer.com/modules/secrets"
	"deployer.com/modules/servers"
	"deployer.com/modules/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func NewFiber() *fiber.App {
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	return app
}

// NewDockerCommunication создает новый экземпляр Docker клиента для DI
func NewDockerCommunication() (*libs.DockerComunication, error) {
	docker, err := libs.NewDockerCommunication()
	if err != nil {
		return nil, err
	}

	// Настройка кэша для deployment-worker контейнеров
	docker.SetCacheExpiration(1 * time.Minute) // Кэш на 1 минуту

	return docker, nil
}

func RegisterRoutes(app *fiber.App, db *gorm.DB, docker *libs.DockerComunication) {
	api := app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Create user service instance for API key validation
	userService := users.NewUsersService(db)

	{
		group := api.Group("/auth")
		routes := auth.NewAuthController(auth.NewAuthService(userService))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/secrets")
		routes := secrets.NewSecretsController(&group, secrets.NewSecretsService(db))
		routes.RegisterRoutes(&group)
	}
	// API key only routes for external integrations
	{
		group := api.Group("/api-secrets")
		routes := secrets.NewSecretsController(&group, secrets.NewSecretsService(db))
		routes.RegisterApiKeyRoutes(&group, userService)
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
		// TODO: В будущем можно передать Docker клиент в контроллер контейнеров
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
		// TODO: В будущем можно передать Docker клиент в контроллер развертываний
		routes := deployments.NewDeploymentsController(&group, deployments.NewDeploymentsService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/projects")
		routes := projects.NewProjectsController(&group, projects.NewProjectsService(db))
		routes.RegisterRoutes(&group)
	}
	{
		group := api.Group("/execute")
		routes := execute.NewExecuteController(
			&group,
			execute.NewExecuteService(scripts.NewScriptsService(db),
				servers.NewServersService(db),
				secrets.NewSecretsService(db),
				containers.NewContainersService(db),
				deployments.NewDeploymentsService(db),
				projects.NewProjectsService(db),
				docker,
			),
		)
		routes.RegisterExecuteRoutes(&group)
	}
}

func main() {
	app := fx.New(
		fx.Provide(
			NewFiber,
			postgres.NewGormDB,
			NewDockerCommunication, // Добавляем Docker клиент в DI контейнер
		),
		fx.Invoke(func(lc fx.Lifecycle, app *fiber.App, db *gorm.DB, docker *libs.DockerComunication) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					// Автомиграция базы данных
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

					// Инициализация кэша deployment-worker контейнеров
					log.Println("Initializing deployment-worker containers cache...")
					if err := docker.RefreshDeploymentWorkersCache(ctx); err != nil {
						log.Printf("Warning: Failed to initialize deployment-worker cache: %v", err)
					} else {
						log.Println("Deployment-worker cache initialized successfully")

						// Получаем и выводим все найденные контейнеры
						containers := docker.GetCachedDeploymentWorkers()
						if len(containers) == 0 {
							log.Println("No deployment-worker containers found")
						} else {
							log.Printf("Found %d deployment-worker container(s):", len(containers))
							for i, container := range containers {
								log.Printf("  [%d] Name: %s | ID: %s | Status: %s | Created: %s",
									i+1,
									container.Name,
									container.ID,
									container.Status,
									container.CreatedAt)
							}
						}
					}

					// Запуск автоматического обновления кэша каждые 30 секунд
					docker.StartAutoRefresh(ctx, 30*time.Second)
					log.Println("Docker cache auto-refresh started (30s interval)")

					// Регистрация маршрутов
					RegisterRoutes(app, db, docker)

					// Запуск веб-сервера
					go func() {
						log.Println("Starting server on :8080...")
						if err := app.Listen(":8080"); err != nil {
							log.Fatal(err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Println("Shutting down services...")

					// Закрытие Docker клиента
					if err := docker.Close(); err != nil {
						log.Printf("Error closing Docker client: %v", err)
					}

					// Остановка веб-сервера
					return app.Shutdown()
				},
			})
		}),
	)
	app.Run()
}
