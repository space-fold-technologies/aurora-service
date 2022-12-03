package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
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
	pluginRegistry  plugins.PluginRegistry
	provider        providers.Provider
	sessionDuration time.Duration
	agent           providers.AgentClient
	scheduler       *gocron.Scheduler
	appservice      *apps.AppService
}

func ProduceServiceResources(
	server *server.ServerCore,
	parameters configuration.Configuration,
	tokenHandler security.TokenHandler,
	hasher security.HashHandler,
	pluginRegistry plugins.PluginRegistry) *ServiceResources {
	return &ServiceResources{
		server:          server,
		parameters:      parameters,
		hasher:          hasher,
		tokenHandler:    tokenHandler,
		passwordHandler: security.NewPasswordHandler(parameters.EncryptionParameters),
		pluginRegistry:  pluginRegistry,
	}
}

func (sr *ServiceResources) Initialize() {
	sr.sessionDuration = time.Duration(sr.parameters.SessionDuration) * time.Hour
	sr.agent = providers.NewClient(sr.parameters.AgentParameters)
	sr.dataSource = sr.createDataSource()
	sr.provider = sr.providers(sr.parameters.Provider)
	if sr.hasPlugin("traefik") {
		sr.pluginRegistry.Put(plugins.REVERSE_PROXY, plugins.NewTraefikPlugin(
			fmt.Sprintf("%s-%s", sr.parameters.NetworkPrefix, sr.parameters.NetworkName),
			sr.parameters.Domain,
			sr.parameters.Https,
			sr.parameters.CertResolverName,
			sr.parameters.CertResolverEmail,
			sr.provider))
	}
	sr.appservice = apps.NewService(sr.provider, sr.hasher, apps.NewRepository(sr.dataSource))
	sr.setupControllers(sr.server.GetRegistry())
	sr.scheduler = gocron.NewScheduler(time.UTC)
	sr.scheduler.Every(5).Minutes().Do(sr.appservice.UpdateContainers)
	sr.scheduler.StartAsync()
}

func (sr *ServiceResources) providers(name string) providers.Provider {
	if name == "DOCKER-SWARM" {
		return docker.NewProvider(sr.pluginRegistry, sr.agent, docker.DockerServiceConfigurations{
			NetworkName:        sr.parameters.NetworkName,
			NetworkPrefix:      sr.parameters.NetworkPrefix,
			ListenerAddress:    sr.parameters.DefaultClusterListenAddress,
			AdvertisingAddress: sr.parameters.DefaultClusterAdvertiseAddress,
		})
	}
	return nil
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
	registry.AddController(apps.NewController(sr.appservice))
	registry.AddController(clusters.NewController(clusters.NewService(sr.provider, clusters.NewRepository(sr.dataSource))))
	registry.AddController(environments.NewController(environments.NewService(environments.NewRepository(sr.dataSource))))
	registry.AddController(nodes.NewController(nodes.NewService(sr.provider, nodes.NewRepository(sr.dataSource))))
	registry.AddController(teams.NewController(teams.NewService(teams.NewRepository(sr.dataSource))))
	registry.AddController(users.NewController(users.NewService(users.NewRepository(sr.dataSource), sr.passwordHandler)))
	registry.AddController(authorization.NewController(authorization.NewService(
		authorization.NewRepository(sr.dataSource),
		sr.passwordHandler,
		sr.tokenHandler,
		sr.sessionDuration,
	)))
}

func (sr *ServiceResources) hasPlugin(name string) bool {
	for _, plugin := range sr.parameters.Plugins {
		if strings.Trim(plugin, " ") == name {
			return true
		}
	}
	return false
}

func (sr *ServiceResources) OnShutDown() {
	sr.scheduler.Stop()
}
