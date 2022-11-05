package app

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
	configurationPath := filepath.Join(i.configs.ProfileDIR, "configurations")
	if _, err := os.Stat(configurationPath); os.IsNotExist(err) {
		if err := os.MkdirAll(configurationPath, fs.ModeAppend); err != nil {
			logging.GetInstance().Error(err)
			os.Exit(-1)
		}
	}
	return i
}

func (i *Initializer) Migrate() {
	dbFile := filepath.Join(i.configs.ProfileDIR, "configurations", "store.db")
	if err := i.migrator.Migrate(dbFile); err != nil && !strings.Contains(err.Error(), "no change") {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else {
		datasource := i.dataSource()
		if has, err := i.hasPersmissions(datasource); err != nil {
			logging.GetInstance().Error(err)
			os.Exit(-1)
		} else if !has {
			if err := i.permissions(i.dataSource()); err != nil {
				logging.GetInstance().Error(err)
				os.Exit(-1)
			} else {
				logging.GetInstance().Infof("Added permissions")
			}
		}
	}
}

func (i *Initializer) Reset() {
	dbFile := filepath.Join(i.configs.ProfileDIR, "configurations", "store.db")
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
	if err := teams.NewRepository(datasource).Create(&teams.Team{
		Name:        i.configs.DefaultTeam,
		Identifier:  uuid.NewString(),
		Description: "default team with full control"}); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
	if err := i.registerAdmin(email, password, userService); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	} else if err := i.setupCluster(clusterService, nodeService, i.configs.DefaultTeam); err != nil {
		logging.GetInstance().Error(err)
		os.Exit(-1)
	}
}

func (i *Initializer) setupCluster(clusterService *clusters.ClusterService, nodeService *nodes.NodeService, team string) error {
	if err := clusterService.Create(&clusters.CreateClusterOrder{
		Name:        i.configs.DefaultCluster,
		Description: "Default cluster created for deployments",
		Teams:       []string{i.configs.DefaultTeam},
		Type:        clusters.CreateClusterOrder_DOCKER_SWARM,
		Namespace:   "swarm",
		Address:     i.configs.DefaultClusterAdvertiseAddress,
	}); err != nil {
		logging.GetInstance().Info("initialization fools")
		return err
	} else if err := nodeService.Create(&nodes.CreateNodeOrder{
		Name:        "default",
		Type:        "DOCKER",
		Description: "Default node",
		Cluster:     i.configs.DefaultCluster,
		Address:     i.configs.DefaultClusterAdvertiseAddress,
	}); err != nil {
		return err
	}
	return nil
}

func (i *Initializer) registerAdmin(email, password string, service *users.UserService) error {
	list := make([]string, 0)
	if err := service.Create(&users.CreateUserOrder{
		Name:     "administrator",
		Nickname: "admin",
		Email:    email,
		Password: password}); err != nil {
		return err
	} else if data, err := os.ReadFile(filepath.Join(i.configs.ProfileDIR, "configurations", "permissions.json")); err != nil {
		return err
	} else if err := json.Unmarshal(data, &list); err != nil {
		return err
	} else if err := service.AddPermissions(&users.UpdatePermissions{Email: email, Permissions: list}); err != nil {
		return err
	} else if err := service.AddTeams(&users.UpdateTeams{Email: email, Teams: []string{i.configs.DefaultTeam}}); err != nil {
		return err
	}
	return nil
}

func (i *Initializer) provider(name string) providers.Provider {
	if name == "DOCKER-SWARM" {
		logging.GetInstance().Infof("Docker Swarm Provider")
		return docker.NewProvider(i.plugin())
	}
	logging.GetInstance().Infof("No Supported Provider found")
	return nil
}

func (i *Initializer) plugin() providers.PluginParameterInjector {
	instance := new(plugins.TraefikPlugin)
	instance.Https = i.configs.Https
	instance.CertResolverName = i.configs.CertResolver
	return instance
}

func (i *Initializer) dataSource() database.DataSource {
	dbFile := filepath.Join(i.configs.ProfileDIR, "configurations", "store.db")
	return database.NewBuilder().EnableLogging().Path(dbFile).Build()
}

func (i *Initializer) permissions(datasource database.DataSource) error {
	entries := []string{}
	arguments := []interface{}{}
	list := make([]string, 0)
	if data, err := os.ReadFile(filepath.Join(i.configs.ProfileDIR, "configurations", "permissions.json")); err != nil {
		return err
	} else if err := json.Unmarshal(data, &list); err != nil {
		return err
	} else {
		for _, entry := range list {
			entries = append(entries, "(?)")
			arguments = append(arguments, entry)
		}
		var sql = "INSERT INTO permission_tb(name) VALUES %s"
		sql = fmt.Sprintf(sql, strings.Join(entries, ","))
		tx := datasource.Connection().Begin()
		if err := tx.Exec(sql, arguments...).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}
}

func (i *Initializer) hasPersmissions(datasource database.DataSource) (bool, error) {
	sql := "EXISTS(SELECT 1 FROM permission_tb) AS has"
	connection := datasource.Connection()
	var has bool
	if err := connection.Raw(sql, has).Error; err != nil {
		return false, err
	}
	return has, nil
}
