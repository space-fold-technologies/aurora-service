package clusters_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/domain/clusters"
	"github.com/space-fold-technologies/aurora-service/app/domain/nodes"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestClusterRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster Repository Test Suite")
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
	if err := connection.AutoMigrate(&teams.Team{}, &clusters.Cluster{}, &nodes.Node{}, &nodes.Container{}); err != nil {
		logging.GetInstance().Info("Crap")
		logging.GetInstance().Error(err)
	} else {
		logging.GetInstance().Info("MIgrations complete")
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

func addTeam(Name string, datasource database.DataSource) {
	tx := datasource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logging.GetInstance().Error(err)
		return
	} else if err := tx.Create(&teams.Team{Identifier: uuid.NewString(), Name: Name, Description: "sample team"}).Error; err != nil {
		logging.GetInstance().Error(err)
		tx.Rollback()
	} else if err := tx.Commit().Error; err != nil {
		logging.GetInstance().Error(err)
	} else {
		logging.GetInstance().Info("Team Added")
	}
}

func addNode(Name, IP, Cluster string, datasource database.DataSource) {
	sql := "INSERT INTO node_tb(name, identifier, description, type, address, cluster_id) " +
		"VALUES(?, ?, ?, ?, ?, (SELECT id FROM cluster_tb WHERE name = ?))"

	tx := datasource.Connection().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Exec(sql, Name, uuid.NewString(), "Sample Node", "docker", "http://10.45.7.12:2773", Cluster).Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error(err)
	} else if err := tx.Commit().Error; err != nil {
		logging.GetInstance().Error(err)
	}
}

var _ = Describe("Cluster store tests", func() {
	dataSource := NewDataSource()
	migrate(dataSource)
	repository := clusters.NewRepository(dataSource)
	Context("Given no existing cluster", func() {
		BeforeEach(func() {
			addTeam("sigma", dataSource)
		})

		AfterEach(func() {
			reset(dataSource)
		})

		It("A Cluster can be added", func() {
			err := repository.Create(&clusters.ClusterEntry{
				Identifier:  uuid.New().String(),
				Name:        "Katsura",
				Description: "Team Katsura Corp Cluster",
				Type:        "DOCKER-SWARM",
				Address:     "http://10.45.9.80",
				Namespace:   "swm",
				Teams:       []string{"sigma"},
			})
			Expect(err).ToNot(HaveOccurred())
			result, err := repository.List([]string{"sigma"})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(1))
			Expect(result[0].Name).To(BeEquivalentTo("Katsura"))
			Expect(result[0].Address).To(BeEquivalentTo("http://10.45.9.80"))
			Expect(result[0].Description).To(BeEquivalentTo("Team Katsura Corp Cluster"))
			Expect(result[0].Type).To(BeEquivalentTo("DOCKER-SWARM"))
			Expect(result[0].Namespace).To(BeEquivalentTo("swm"))
			Expect(len(result[0].Teams)).To(BeEquivalentTo(1))
		})
	})
	Context("Given an existing cluster", func() {

		BeforeEach(func() {
			addTeam("sigma", dataSource)
			repository.Create(&clusters.ClusterEntry{
				Identifier:  uuid.New().String(),
				Name:        "Katsura",
				Description: "Team Katsura Corp Cluster",
				Type:        "DOCKER-SWARM",
				Address:     "http://10.45.9.80",
				Namespace:   "swm",
				Teams:       []string{"sigma"},
			})
		})
		AfterEach(func() {
			reset(dataSource)
		})

		It("Cannot create another cluster with the same name", func() {
			err := repository.Create(&clusters.ClusterEntry{
				Identifier:  uuid.New().String(),
				Name:        "Katsura",
				Description: "Team Katsura Corp Cluster",
				Type:        "DOCKER-SWARM",
				Address:     "http://10.45.9.88",
				Namespace:   "swm",
				Teams:       []string{"sigma"},
			})
			Expect(err).To(HaveOccurred())
		})
		It("Can update an existing cluster", func() {
			err := repository.Update(
				"Katsura",
				"Katsura Corp, Inc",
				[]string{},
			)
			Expect(err).ToNot(HaveOccurred())
			result, err := repository.List([]string{"sigma"})
			Expect(err).ToNot(HaveOccurred())
			Expect(result[0].Name).To(BeEquivalentTo("Katsura"))
			Expect(result[0].Address).To(BeEquivalentTo("http://10.45.9.80"))
			Expect(result[0].Description).To(BeEquivalentTo("Katsura Corp, Inc"))
			Expect(result[0].Type).To(BeEquivalentTo("DOCKER-SWARM"))
			Expect(result[0].Namespace).To(BeEquivalentTo("swm"))
			Expect(len(result[0].Teams)).To(BeEquivalentTo(1))
		})
		It("Can remove a cluster that has no nodes", func() {
			err := repository.Remove("Katsura")
			Expect(err).ToNot(HaveOccurred())
			result, err := repository.List([]string{"sigma"})
			Expect(err).ToNot(HaveOccurred())
			Expect(len(result)).To(BeEquivalentTo(0))
		})
		It("Can check if a cluster has nodes", func() {
			addNode("Voltmax", "10.45.90.12", "Katsura", dataSource)
			hasNode, err := repository.HasNodes("Katsura")
			Expect(err).ToNot(HaveOccurred())
			Expect(hasNode).To(BeTrue())
		})
	})
})
