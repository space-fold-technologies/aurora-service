package nodes_test

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
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestApplicationRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nodes Repository Test Suite")
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
	connection := dataSource.Connection()
	if err := connection.AutoMigrate(&teams.Team{}, &apps.Application{}, &apps.Deployment{}, &clusters.Cluster{}, &nodes.Node{}, &nodes.Container{}); err != nil {
		logging.GetInstance().Error(err)
	} else {
		logging.GetInstance().Info("MIgrations complete")
	}
}

func addTeam(Name string, datasource database.DataSource) {
	if err := teams.NewRepository(datasource).Create(&teams.Team{Identifier: uuid.NewString(), Name: Name, Description: "sample team"}); err != nil {
		logging.GetInstance().Error(err)
	}
}

func addCluster(Name, Team string, datasource database.DataSource) {
	if err := clusters.NewRepository(datasource).Create(&clusters.ClusterEntry{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Woof Woof",
		Type:        "DOCKER-SWARM",
		Address:     "http://10.45.9.72",
		Namespace:   "swm",
		Teams:       []string{Team},
	}); err != nil {
		logging.GetInstance().Error(err)
	}
}

func addApplication(Name, Team, Cluster string, Scale int32, datasource database.DataSource) {
	apps.NewRepository(datasource).Create(&apps.ApplicationEntry{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Sample Application",
		Team:        Team,
		Cluster:     Cluster,
		Scale:       uint32(Scale),
	})
}

func addContainer(ServiceID, Node, IP string, datasource database.DataSource) {
	apps.NewRepository(datasource).AddContainers([]*apps.ContainerOrder{{
		Identifier: uuid.NewString(),
		IP:         IP,
		Family:     4,
		ServiceID:  ServiceID,
		Node:       Node,
	}})
}

func addDeployment(Application, ServiceID string, datasource database.DataSource) {
	sql := "INSERT INTO deployment_tb(identifier, application_id, image_uri, added_at, completed_at, status, report, service_identifier) " +
		"VALUES(?, (SELECT id FROM application_tb WHERE name = ?), ?, ?, ?, ?, ?, ?)"
	tx := datasource.Connection().Begin()
	addedAt := time.Now()
	completedAt := addedAt.Add(2 * time.Minute)
	if err := tx.Exec(sql, uuid.NewString(), Application, "test@5859303853", addedAt, completedAt, "DEPLOYED", "service deployed successfully", ServiceID).Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error(err)
	} else if err := tx.Commit().Error; err != nil {
		logging.GetInstance().Error(err)
	}
}

func reset(datasource database.DataSource) {
	tx := datasource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			logging.GetInstance().Error("FAILED TO RESET ON PANIC")
			tx.Rollback()
		}
	}()
	if err := tx.Exec("DELETE FROM container_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED CONTAINER DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM deployment_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED DEPLOYMENT DELETE")
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

var _ = Describe("Team store tests", func() {
	dataSource := NewDataSource()
	repository := nodes.NewRepository(dataSource)
	migrate(dataSource)
	Context("Given no existing nodes", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		BeforeEach(func() {
			addTeam("Bakumatsu", dataSource)
			addCluster("Ryogen", "Bakumatsu", dataSource)
		})

		It("Can add node to cluster", func() {
			err := repository.Create(&nodes.NodeEntry{
				Identifier:  uuid.NewString(),
				Name:        "node-01",
				Description: "First Node",
				Type:        "Docker",
				Address:     "http://10.45.9.23:2773",
				Cluster:     "Ryogen",
			})

			Expect(err).ToNot(HaveOccurred())
			results, err := repository.List("Ryogen")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(results)).To(BeEquivalentTo(1))
			Expect(results[0].Name).To(BeEquivalentTo("node-01"))
			Expect(results[0].Address).To(BeEquivalentTo("http://10.45.9.23:2773"))
			Expect(results[0].Description).To(BeEquivalentTo("First Node"))
			Expect(results[0].Type).To(BeEquivalentTo("Docker"))
		})
	})
	Context("Given an existing node", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		BeforeEach(func() {
			addTeam("Bakumatsu", dataSource)
			addCluster("Ryogen", "Bakumatsu", dataSource)
			repository.Create(&nodes.NodeEntry{
				Name:        "node-01",
				Description: "First Node",
				Type:        "Docker",
				Address:     "http://10.45.9.23:2773",
				Cluster:     "Ryogen",
			})
		})

		It("Cannot add a node with the same name", func() {
			err := repository.Create(&nodes.NodeEntry{
				Name:        "node-01",
				Description: "Testing limits",
				Type:        "Podman",
				Address:     "http://10.45.9.24:2773",
				Cluster:     "Ryogen",
			})
			Expect(err).To(HaveOccurred())
		})
		It("Cannot add a node with the same address", func() {
			err := repository.Create(&nodes.NodeEntry{
				Name:        "node-02",
				Description: "Testing limits",
				Type:        "Podman",
				Address:     "http://10.45.9.23:2773",
				Cluster:     "Ryogen",
			})
			Expect(err).To(HaveOccurred())
		})
		It("Can update node", func() {
			err := repository.Update("node-01", "Fire Fire Fire")
			Expect(err).ToNot(HaveOccurred())
			results, err := repository.List("Ryogen")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(results)).To(BeEquivalentTo(1))
			Expect(results[0].Name).To(BeEquivalentTo("node-01"))
			Expect(results[0].Address).To(BeEquivalentTo("http://10.45.9.23:2773"))
			Expect(results[0].Description).To(BeEquivalentTo("Fire Fire Fire"))
			Expect(results[0].Type).To(BeEquivalentTo("Docker"))
		})
		It("Can remove node", func() {
			err := repository.Remove("node-01")
			Expect(err).ToNot(HaveOccurred())
			results, err := repository.List("Ryogen")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(results)).To(BeEquivalentTo(0))
		})
		It("Can check for containers", func() {
			addApplication("subaru", "Bakumatsu", "Ryogen", 1, dataSource)
			serviceId := uuid.NewString()
			addDeployment("subaru", serviceId, dataSource)
			addContainer(serviceId, "node-01", "172.168.44.90", dataSource)
			hasContainers, _ := repository.HasContainers("node-01")
			Expect(hasContainers).To(BeTrue())
		})
	})
})
