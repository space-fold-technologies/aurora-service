package teams_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTeamRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Team Repository Test Suite")
}

type testDataSource struct {
	db *gorm.DB
}

func (tds *testDataSource) Connection() *gorm.DB {
	if tds.db == nil {
		var err error
		tds.db, err = gorm.Open(
			sqlite.Open("file::memory:?mode=memory&cache=shared"),
			&gorm.Config{})
		if err != nil {
			panic(err)
		}
	}

	return tds.db
}

func (tds *testDataSource) Close() error {
	if tds.db != nil {
		return tds.Close()
	}
	return nil
}

func dataSourceProvider() database.DataSource {
	return &testDataSource{}
}

func migrate(dataSource database.DataSource) {
	// TODO: inline migration scripts
	dataSource.Connection().AutoMigrate(&teams.Team{})
}

var _ = Describe("Team store tests", func() {
	dataSource := dataSourceProvider()
	migrate(dataSource)
	repository := teams.NewRepository(dataSource)
	Context("Given no existing team", func() {
		AfterEach(func() {
			dataSource.Connection().Exec("DELETE FROM team_tb")
		})
		It("A team can be added", func() {
			err := repository.Create(&teams.Team{
				Name:        "orion",
				Description: "Communiation Technology Development Department",
			})
			Expect(err).ToNot(HaveOccurred())
			teams, err := repository.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(teams).ToNot(BeEmpty())
			Expect(teams[0].Name).To(Equal("orion"))
			Expect(teams[0].Description).To(Equal("Communiation Technology Development Department"))
		})
		It("A team can be removed", func() {
			err := repository.Create(&teams.Team{
				Name:        "orion",
				Description: "Communiation Technology Development Department",
			})
			Expect(err).ToNot(HaveOccurred())
			err = repository.Remove("orion")
			Expect(err).ToNot(HaveOccurred())
			teams, err := repository.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(teams).To(BeEmpty())
		})
	})
	Context("Given an existing team", func() {
		AfterEach(func() {
			dataSource.Connection().Exec("DELETE FROM team_tb")
		})
		BeforeEach(func() {
			repository.Create(&teams.Team{
				Name:        "orion",
				Description: "Communiation Technology Development Department",
			})
		})

		It("A team with the same name cannot be added", func() {
			err := repository.Create(&teams.Team{
				Name:        "orion",
				Description: "Communiation Technology Development Department",
			})
			Expect(err).To(HaveOccurred())
		})
		It("Another unique team can be added", func() {
			err := repository.Create(&teams.Team{
				Name:        "vectron",
				Description: "monitoring tools and data visualizations Department",
			})
			Expect(err).ToNot(HaveOccurred())
			teams, err := repository.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(teams).ToNot(BeEmpty())
			Expect(len(teams)).To(BeNumerically("==", 2))
		})
	})
})
