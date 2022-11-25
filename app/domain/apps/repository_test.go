package apps_test

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

func TestApplicationRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Application Repository Test Suite")
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
	dataSource.Connection().AutoMigrate(&teams.Team{}, &apps.Application{}, &apps.Deployment{}, &clusters.Cluster{}, &nodes.Node{}, &nodes.Container{}, &environments.EnvVar{})
	rule := `CREATE TRIGGER update_last_deployment AFTER UPDATE ON deployment_tb
			  BEGIN
			   UPDATE application_tb SET last_deployment=(CASE WHEN NEW.status = 'DEPLOYED' THEN NEW.completed_at ELSE last_deployment END) WHERE id=NEW.application_id;
			  END;`
	dataSource.Connection().Exec(rule)
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
	} else if err := tx.Exec("DELETE FROM application_tb").Error; err != nil {
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

func addTeam(datasource database.DataSource, Name string) {
	teams.NewRepository(datasource).Create(&teams.Team{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Top Shelf Team",
	})
}

func addCluster(datasource database.DataSource, Name, Team string) {
	clusters.NewRepository(datasource).Create(&clusters.ClusterEntry{
		Identifier:  uuid.NewString(),
		Name:        Name,
		Description: "Top Shelf Cluster",
		Type:        "DOCKER-SWARM",
		Teams:       []string{Team},
	})
}

func addNode(datasource database.DataSource, Name, Identifier, Cluster string) {

	nodes.NewRepository(datasource).Create(&nodes.NodeEntry{
		Identifier:  Identifier,
		Name:        Name,
		Description: "Top Shelf Node",
		Type:        "DOCKER",
		Address:     "10.45.90.8",
		Cluster:     Cluster,
	})
}

func addEnvironment(datasource database.DataSource, Scope, Target, Key, Value string) {
	environments.NewRepository(datasource).Add(Scope, Target, []*environments.Entry{{Key: Key, Value: Value}})
}

var _ = Describe("Application store tests", func() {
	dataSource := NewDataSource()
	migrate(dataSource)
	repository := apps.NewRepository(dataSource)
	Context("Given no existing application", func() {
		BeforeEach(func() {
			addTeam(dataSource, "axolotol")
			addCluster(dataSource, "pond", "axolotol")
		})

		AfterEach(func() {
			reset(dataSource)
		})
		It("An application entry can be added", func() {
			identifier := uuid.NewString()
			err := repository.Create(&apps.ApplicationEntry{
				Identifier:  identifier,
				Name:        "dio",
				Description: "Sample application profile",
				Team:        "axolotol",
				Cluster:     "pond",
				Scale:       2,
			})
			Expect(err).ToNot(HaveOccurred())
			details, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(details.Name).To(BeEquivalentTo("dio"))
			Expect(details.Description).To(BeEquivalentTo("Sample application profile"))
			Expect(details.Team).ToNot(BeNil())
			Expect(details.Cluster).ToNot(BeNil())
			Expect(details.Identifier).To(BeEquivalentTo(identifier))
			Expect(details.Scale).To(BeEquivalentTo(2))
		})
	})
	Context("Given an existing application", func() {
		nodeIdentifier := uuid.NewString()
		BeforeEach(func() {
			addTeam(dataSource, "axolotol")
			addCluster(dataSource, "pond", "axolotol")
			addNode(dataSource, "x1", nodeIdentifier, "pond")
			repository.Create(&apps.ApplicationEntry{
				Identifier:  uuid.NewString(),
				Name:        "dio",
				Description: "Sample application profile",
				Team:        "axolotol",
				Cluster:     "pond",
				Scale:       2,
			})
		})
		AfterEach(func() {
			reset(dataSource)
		})
		It("An application can be updated", func() {
			err := repository.Update("dio", "flex-master", 4)
			Expect(err).ToNot(HaveOccurred())
			details, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(details.Name).To(BeEquivalentTo("dio"))
			Expect(details.Description).To(BeEquivalentTo("flex-master"))
			Expect(details.Team).ToNot(BeNil())
			Expect(details.Cluster).ToNot(BeNil())
			Expect(details.Scale).To(BeEquivalentTo(4))
		})
		It("An application can have containers", func() {
			currentTime := time.Now()
			identifier := uuid.NewString()
			serviceId := uuid.NewString()
			repository.RegisterDeploymentEntry(&apps.DeploymentOrder{
				Identifier:      identifier,
				ApplicationName: "dio",
				Status:          "PENDING",
				ImageURI:        "sme:latest",
				CreatedAt:       &currentTime,
			})
			updatedAt := time.Now()
			repository.UpdateDeploymentEntry(&apps.DeploymentUpdate{
				Identifier: identifier,
				Status:     "DEPLOYED",
				Report:     "Application successfully deployed",
				ServiceID:  serviceId,
				UpdatedAt:  &updatedAt,
			})
			err := repository.AddContainers([]*apps.ContainerOrder{
				{
					Identifier: uuid.NewString(),
					IP:         "176.22.45.90",
					Family:     4,
					ServiceID:  serviceId,
					Node:       nodeIdentifier,
				}})
			Expect(err).ToNot(HaveOccurred())
			details, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(details.Instances)).To(BeEquivalentTo(1))
			Expect(details.Instances[0].Family).To(BeEquivalentTo(4))
			Expect(details.Instances[0].IP).To(BeEquivalentTo("176.22.45.90"))
		})
		It("Applications can be listed", func() {
			summaries, err := repository.List("pond")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(summaries)).To(BeEquivalentTo(1))
			Expect(summaries[0].Name).To(BeEquivalentTo("dio"))
			Expect(summaries[0].Scale).To(BeEquivalentTo(2))
		})
		It("An application can be removed", func() {
			err := repository.Remove("dio")
			Expect(err).ToNot(HaveOccurred())
			_, err = repository.FetchDetails("dio")
			Expect(err).To(HaveOccurred())
		})
		It("Can register an application deployment", func() {
			identifier := uuid.NewString()
			currentTime := time.Now()
			err := repository.RegisterDeploymentEntry(&apps.DeploymentOrder{
				Identifier:      identifier,
				ApplicationName: "dio",
				Status:          "PENDING",
				ImageURI:        "sme:latest",
				CreatedAt:       &currentTime,
			})
			Expect(err).ToNot(HaveOccurred())
		})
	})
	Context("Given an application deployment has been registered", func() {
		identifier := uuid.NewString()
		nodeIdentifier := uuid.NewString()
		BeforeEach(func() {
			addTeam(dataSource, "axolotol")
			addCluster(dataSource, "pond", "axolotol")
			addNode(dataSource, "x1", nodeIdentifier, "pond")
			repository.Create(&apps.ApplicationEntry{
				Identifier:  uuid.NewString(),
				Name:        "dio",
				Description: "Sample application profile",
				Team:        "axolotol",
				Cluster:     "pond",
				Scale:       2,
			})
			currentTime := time.Now()
			repository.RegisterDeploymentEntry(&apps.DeploymentOrder{
				Identifier:      identifier,
				ApplicationName: "dio",
				Status:          "PENDING",
				ImageURI:        "sme:latest",
				CreatedAt:       &currentTime,
			})
		})
		AfterEach(func() {
			reset(dataSource)
		})
		It("Can update an applcation deployment", func() {
			updatedAt := time.Now()
			err := repository.UpdateDeploymentEntry(&apps.DeploymentUpdate{
				Identifier: identifier,
				Status:     "DEPLOYED",
				Report:     "Application successfully deployed",
				UpdatedAt:  &updatedAt,
			})
			Expect(err).ToNot(HaveOccurred())
			detail, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(detail.LastDeployment.Unix()).To(BeEquivalentTo(updatedAt.Unix()))
		})
		It("Can fetch the image uri for an application", func() {
			uri, err := repository.FetchImageURI(identifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(uri).To(BeEquivalentTo("sme:latest"))
		})
		It("Can fetch an applications environment variables", func() {
			addEnvironment(dataSource, "CLUSTER", "pond", "CONSUL_PORT", "8500")
			addEnvironment(dataSource, "NODE", "x1", "CONSUL_HOST", "10.45.34.90")
			addEnvironment(dataSource, "APPLICATION", "dio", "NAME", "zoku")
			vars, err := repository.FetchEnvVars("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(vars).ToNot(BeNil())
			Expect(len(vars)).To(BeEquivalentTo(3))
		})
	})
	Context("Given an application entry has been made", func() {
		identifier := uuid.NewString()
		serviceId := uuid.NewString()
		nodeIdentifier := uuid.NewString()
		BeforeEach(func() {
			addTeam(dataSource, "axolotol")
			addCluster(dataSource, "pond", "axolotol")
			addNode(dataSource, "x1", nodeIdentifier, "pond")
			repository.Create(&apps.ApplicationEntry{
				Identifier:  uuid.NewString(),
				Name:        "dio",
				Description: "Sample application profile",
				Team:        "axolotol",
				Cluster:     "pond",
				Scale:       2,
			})

			currentTime := time.Now()
			repository.RegisterDeploymentEntry(&apps.DeploymentOrder{
				Identifier:      identifier,
				ApplicationName: "dio",
				Status:          "PENDING",
				ImageURI:        "sme:latest",
				CreatedAt:       &currentTime,
			})
			updatedAt := time.Now()
			repository.UpdateDeploymentEntry(&apps.DeploymentUpdate{
				Identifier: identifier,
				Status:     "DEPLOYED",
				Report:     "Application successfully deployed",
				ServiceID:  serviceId,
				ImageURI:   "sme:digest",
				UpdatedAt:  &updatedAt,
			})
		})
		AfterEach(func() {
			reset(dataSource)
		})
		It("Can add a container entry for an application", func() {
			err := repository.AddContainers([]*apps.ContainerOrder{
				{
					Identifier: uuid.NewString(),
					IP:         "176.90.87.4",
					Family:     4,
					ServiceID:  serviceId,
					Node:       nodeIdentifier,
				}})
			Expect(err).ToNot(HaveOccurred())
			details, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(len(details.Instances)).To(BeEquivalentTo(1))
			Expect(details.Instances[0].Family).To(BeEquivalentTo(4))
			Expect(details.Instances[0].IP).To(BeEquivalentTo("176.90.87.4"))
		})
		It("Can remove an application container entry", func() {
			containerIdentifier := uuid.NewString()
			err := repository.AddContainers([]*apps.ContainerOrder{{
				Identifier: containerIdentifier,
				IP:         "176.90.87.4",
				Family:     4,
				ServiceID:  serviceId,
				Node:       nodeIdentifier,
			}})
			Expect(err).ToNot(HaveOccurred())
			err = repository.RemoveContainer(containerIdentifier)
			Expect(err).ToNot(HaveOccurred())
			details, err := repository.FetchDetails("dio")
			Expect(err).ToNot(HaveOccurred())
			Expect(details.Instances).To(BeEmpty())
		})
	})
})
