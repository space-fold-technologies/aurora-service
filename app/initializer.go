package app

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/google/uuid"
	"github.com/space-fold-technologies/aurora-service/app/core/configuration"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/core/plugins"
	"github.com/space-fold-technologies/aurora-service/app/core/providers"
	"github.com/space-fold-technologies/aurora-service/app/core/providers/docker"
	"github.com/space-fold-technologies/aurora-service/app/core/security"
	"github.com/space-fold-technologies/aurora-service/app/domain/clusters"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"github.com/space-fold-technologies/aurora-service/app/domain/users"
)

type Initializer struct {
	configs         configuration.Configuration
	passwordHandler security.PasswordHandler
	migrator        database.MigrationHandler
}

func NewInitializer() *Initializer {
	instance := &Initializer{}
	return instance.init()
}

func (i *Initializer) init() *Initializer {
	i.configs = configuration.ParseFromResource()
	i.passwordHandler = security.NewPasswordHandler(i.configs.EncryptionParameters)
	i.migrator = database.NewMigrationHandler()
	configurationPath := fmt.Sprintf("%s/configurations", i.configs.ProfileDIR)
	if _, err := os.Stat(configurationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configurationPath, fs.ModeAppend); err != nil {
			logging.GetInstance().Error(err)
			os.Exit(-1)
		}
	}
	return i
}

func (i *Initializer) Migrate() {
	dbFile := fmt.Sprintf("%s/configurations/store.db", i.configs.ProfileDIR)
	if err := i.migrator.Migrate(dbFile); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
}

func (i *Initializer) Reset() {
	dbFile := fmt.Sprintf("%s/configurations/store.db", i.configs.ProfileDIR)
	if err := i.migrator.Reset(dbFile); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
}

func (i *Initializer) Register(email, password string) {
	// Capture and validate the inputs
	// wire up a store to save credentials
	// create cluster
	// create single server node
	datasource := i.dataSource()
	provider := i.provider(i.configs.Provider)
	userService := users.NewService(users.NewRepository(datasource), i.passwordHandler)
	clusterService := clusters.NewService(provider, clusters.NewRepository(datasource))
	nodeService := nodes.NewService(provider, nodes.NewRepository(datasource))
	if err := i.registerAdmin(email, password, userService); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else if err := i.setupCluster(clusterService, nodeService, i.configs.DefaultTeam); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
}

func (i *Initializer) setupCluster(clusterService *clusters.ClusterService, nodeService *nodes.NodeService, team string) error {
	if err := teams.NewRepository(i.dataSource()).Create(&teams.Team{
		Identifier:  uuid.NewString(),
		Name:        i.configs.DefaultTeam,
		Description: "Default team with full control",
	}); err != nil {
		return err
	}
	if err := clusterService.Create(&clusters.CreateClusterOrder{
		Name:        i.configs.DefaultCluster,
		Description: "Default cluster created for deployments",
		Teams:       []string{i.configs.DefaultTeam},
		Type:        clusters.CreateClusterOrder_DOCKER_SWARM,
		Namespace:   "swarm",
		Address:     i.configs.DefaultClusterAddress,
	}); err != nil {
		return err
	} else if err := nodeService.Create(&nodes.CreateNodeOrder{
		Name:        "default",
		Type:        "DOCKER",
		Description: "Default node",
		Cluster:     i.configs.DefaultCluster,
		Address:     i.configs.DefaultClusterAddress,
	}); err != nil {
		return err
	}
	return nil
}

func (i *Initializer) registerAdmin(email, password string, service *users.UserService) error {
	return service.Create(&users.CreateUserOrder{Name: "admin", Email: email, Password: password})
}

func (i *Initializer) provider(name string) providers.Provider {
	if name == "DOCKER-SWARN" {
		return docker.NewProvider(i.plugin())
	}
	return nil
}

func (i *Initializer) plugin() providers.PluginParameterInjector {
	instance := new(plugins.TraefikPlugin)
	instance.Https = i.configs.Https
	instance.CertResolverName = i.configs.CertResolver
	return instance
}

func (i *Initializer) dataSource() database.DataSource {
	dbFile := fmt.Sprintf("%s/configurations/store.db", i.configs.ProfileDIR)
	return database.NewBuilder().EnableLogging().Path(dbFile).Build()
}
