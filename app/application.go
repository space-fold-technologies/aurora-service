package app

import (
	"path/filepath"

	"github.com/space-fold-technologies/aurora-service/app/core/configuration"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/plugins"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	"github.com/space-fold-technologies/aurora-service/app/core/providers/docker"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/core/server"
	"github.com/space-fold-technologies/aurora-service/app/core/server/http/controllers"
	"github.com/space-fold-technologies/aurora-service/app/domain/apps"
	"github.com/space-fold-technologies/aurora-service/app/domain/authorization"
	"github.com/space-fold-technologies/aurora-service/app/domain/clusters"
	"github.com/space-fold-technologies/aurora-service/app/domain/environments"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"github.com/space-fold-technologies/aurora-service/app/domain/users"
)

type ServiceResources struct {
	server          *server.ServerCore
	dataSource      database.DataSource
	parameters      configuration.Configuration
	passwordHandler security.PasswordHandler
	tokenHandler    security.TokenHandler
	hasher          security.HashHandler
	provider        providers.Provider
}

func ProduceServiceResources(
	server *server.ServerCore,
	parameters configuration.Configuration,
	tokenHandler security.TokenHandler,
	hasher security.HashHandler) *ServiceResources {
	return &ServiceResources{
		server:          server,
		parameters:      parameters,
		hasher:          hasher,
		tokenHandler:    tokenHandler,
		passwordHandler: security.NewPasswordHandler(parameters.EncryptionParameters),
	}
}

func (sr *ServiceResources) Initialize() {
	sr.dataSource = sr.createDataSource()
	sr.provider = sr.providers(sr.parameters.Provider)
	sr.setupControllers(sr.server.GetRegistry())
}

func (sr *ServiceResources) providers(name string) providers.Provider {
	if name == "DOCKER-SWARM" {
		return docker.NewProvider(sr.plugin())
	}
	return nil
}

func (sr *ServiceResources) plugin() providers.PluginParameterInjector {
	instance := new(plugins.TraefikPlugin)
	instance.Https = sr.parameters.Https
	instance.CertResolverName = sr.parameters.CertResolver
	return instance
}

func (sr *ServiceResources) createDataSource() database.DataSource {
	return database.
		NewBuilder().
		EnableLogging().
		Path(filepath.Join(sr.parameters.ProfileDIR, "configurations", "store.db")).
		Build()
}

func (sr *ServiceResources) setupControllers(registry *controllers.HTTPControllerRegistry) {
	//TODO: Register all repositories and inject inject into services and controllers
	registry.AddController(apps.NewController(apps.NewService(sr.provider, sr.hasher, apps.NewRepository(sr.dataSource))))
	registry.AddController(clusters.NewController(clusters.NewService(sr.provider, clusters.NewRepository(sr.dataSource))))
	registry.AddController(environments.NewController(environments.NewService(environments.NewRepository(sr.dataSource))))
	registry.AddController(nodes.NewController(nodes.NewService(sr.provider, nodes.NewRepository(sr.dataSource))))
	registry.AddController(teams.NewController(teams.NewService(teams.NewRepository(sr.dataSource))))
	registry.AddController(users.NewController(users.NewService(users.NewRepository(sr.dataSource), sr.passwordHandler)))
	registry.AddController(authorization.NewController(authorization.NewService(
		authorization.NewRepository(sr.dataSource),
		sr.passwordHandler,
		sr.tokenHandler,
		sr.parameters.SessionDuration,
	)))
}
