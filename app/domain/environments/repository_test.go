package environments_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/domain/apps"
	"github.com/space-fold-technologies/aurora-service/app/domain/clusters"
	"github.com/space-fold-technologies/aurora-service/app/domain/environments"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEnvironmentRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Environment Repository Test Suite")
}

type testDataSource struct {
	db *gorm.DB
}

func NewDataSource() database.DataSource {
	instance := new(testDataSource)
	instance.initialize()
	return instance
}

func (tds *testDataSource) initialize() {
	if tds.db == nil {
		var err error
		tds.db, err = gorm.Open(
			sqlite.Open("file::memory:?mode=memory&cache=shared"),
			&gorm.Config{})
		if err != nil {
			panic(err)
		} else if db, err := tds.db.DB(); err != nil {
			panic(err)
		} else {
			db.SetConnMaxIdleTime(time.Minute * 5)
			db.SetMaxIdleConns(10)
			db.SetMaxOpenConns(100)
			db.SetConnMaxLifetime(time.Hour)
		}
	}
}

func (tds *testDataSource) Connection() *gorm.DB {
	return tds.db
}

func (tds *testDataSource) Close() error {
	if tds.db != nil {
		return tds.Close()
	}
	return nil
}

func migrate(dataSource database.DataSource) {
	// TODO: inline migration scripts
	dataSource.Connection().AutoMigrate(&apps.Application{}, &clusters.Cluster{}, &nodes.Node{}, &environments.EnvVar{})
}

func reset(datasource database.DataSource) {
	tx := datasource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			logging.GetInstance().Error("FAILED TO RESET ON PANIC")
			tx.Rollback()
		}
	}()
	if err := tx.Exec("DELETE FROM application_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED APPLICATION DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM node_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED NODE DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM cluster_teams").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED CLUSTER TEAMS DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM cluster_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED CLUSTER DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM environment_variable_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED ENVIRONMENT VARIABLES DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM team_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED TEAM DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Commit().Error; err != nil {
		logging.GetInstance().Error(err)
	} else {
		logging.GetInstance().Info("RESET")
	}
}

func addCluster(Name, Team string, dataSource database.DataSource) {
	clusters.NewRepository(dataSource).Create(&clusters.ClusterEntry{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Top Shelf Cluster",
		Type:        "DOCKER-SWARM",
		Teams:       []string{Team},
	})
}

func addTeam(Name string, datasource database.DataSource) {
	teams.NewRepository(datasource).Create(&teams.Team{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Top Shelf Team",
	})
}

func addApplication(Name, Cluster, Team string, datasource database.DataSource) {
	apps.NewRepository(datasource).Create(&apps.ApplicationEntry{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Demo Service",
		Cluster:     Cluster,
		Team:        Team,
		Scale:       2,
	})
}

var _ = Describe("Environment store tests", func() {
	dataSource := NewDataSource()
	migrate(dataSource)
	repository := environments.NewRepository(dataSource)
	Context("Given no existing environment variable entries", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		BeforeEach(func() {
			addTeam("sigma", dataSource)
			addCluster("spiral", "sigma", dataSource)
			addApplication("cart-service", "spiral", "sigma", dataSource)
		})

		It("An environment variable can be added", func() {
			vars := []*environments.Entry{
				{Key: "CONSUL_DATACENTER", Value: "ares"},
			}
			err := repository.Add("CLUSTER", "cl-01", vars)
			Expect(err).ToNot(HaveOccurred())
			vars, err = repository.List("CLUSTER", "cl-01")
			Expect(err).ToNot(HaveOccurred())
			Expect(vars).ToNot(BeEmpty())
			Expect(len(vars)).To(BeNumerically("==", 1))
		})

		It("An environment variable can be removed", func() {
			vars := []*environments.Entry{
				{Key: "CONSUL_DATACENTER", Value: "ares"},
			}
			err := repository.Add("CLUSTER", "cl-01", vars)
			Expect(err).ToNot(HaveOccurred())
			err = repository.Remove([]string{"CONSUL_DATACENTER"})
			Expect(err).ToNot(HaveOccurred())
			vars, err = repository.List("CLUSTER", "cl-01")
			Expect(err).ToNot(HaveOccurred())
			Expect(vars).To(BeEmpty())
		})
	})
	Context("Given existing environment variable entries", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		BeforeEach(func() {
			clusterVars := []*environments.Entry{
				{Key: "CONSUL_DATACENTER", Value: "ares"},
				{Key: "CONSUL_TOKEN", Value: "that**some&&fuxys**%^$&"},
			}
			nodeVars := []*environments.Entry{}
			appVars := []*environments.Entry{}
			repository.Add("CLUSTER", "cl-01", clusterVars)
			repository.Add("NODE", "cl-01-node-01", nodeVars)
			repository.Add("APPLICATION", "cl-01-node-01-app-tsumaka", appVars)
		})

		It("An existing environment variable cannot be duplicated", func() {
			err := repository.Add("NODE", "cl-01-node-01", []*environments.Entry{
				{Key: "CONSUL_TOKEN", Value: "*9040328%^$*#()$@)IJAUJRK#("},
			})
			Expect(err).To(HaveOccurred())
		})
	})
	Context("Given existing environment variable entries with related scopes", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		BeforeEach(func() {
			repository.Add("APPLICATION", "cl-01-node-01-app-tsumaka", []*environments.Entry{{Key: "NH", Value: "VC"}})
		})
		It("Will fetch all variables for a given scope", func() {
			vars, err := repository.List("APPLICATION", "cl-01-node-01-app-tsumaka")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(vars)).To(BeEquivalentTo(1))
		})
	})
})
