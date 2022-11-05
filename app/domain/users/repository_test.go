package users_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
	"github.com/space-fold-technologies/aurora-service/app/domain/teams"
	"github.com/space-fold-technologies/aurora-service/app/domain/users"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTeamRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Repository Test Suite")
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

func migrate(dataSource database.DataSource) {
	// TODO: inline migration scripts
	connection := dataSource.Connection()
	if err := connection.AutoMigrate(&teams.Team{}, &users.User{}); err != nil {
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
	if err := tx.Exec("DELETE FROM user_teams").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED USER_TEAMS DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM user_permissions").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED USER PERMISSION DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM user_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED USER DELETE")
		logging.GetInstance().Error(err)
	} else if err := tx.Exec("DELETE FROM permission_tb").Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error("FAILED PERMISSION DELETE")
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
	if err := teams.NewRepository(datasource).Create(&teams.Team{Identifier: uuid.NewString(), Name: Name, Description: "sample team"}); err != nil {
		logging.GetInstance().Error(err)
	}
}

func addPermission(Name string, datasource database.DataSource) {
	sql := "INSERT INTO permission_tb(name) VALUES(?)"
	tx := datasource.Connection().Begin()
	if err := tx.Exec(sql, Name).Error; err != nil {
		tx.Rollback()
		logging.GetInstance().Error(err)
	} else if err := tx.Commit().Error; err != nil {
		logging.GetInstance().Error(err)
	}
}

var _ = Describe("Team store tests", func() {
	dataSource := NewDataSource()
	migrate(dataSource)
	repository := users.NewRepository(dataSource)
	Context("given no users exit", func() {
		AfterEach(func() {
			reset(dataSource)
		})
		It("can add a user", func() {
			identifier := uuid.New()
			err := repository.Create(&users.UserEntry{
				Name:         "Nicole Pedersen",
				NickName:     "NicPed",
				Identifier:   identifier,
				Email:        "nicole.pedersen@mail.no",
				PasswordHash: "($)#)@__$()#)@_@($(@)",
			})
			Expect(err).ToNot(HaveOccurred())
			user, err := repository.FindById(identifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(user.Identifier).To(BeEquivalentTo(identifier))
			Expect(user.Name).To(BeEquivalentTo("Nicole Pedersen"))
			Expect(user.NickName).To(BeEquivalentTo("NicPed"))
			Expect(user.Email).To(BeEquivalentTo("nicole.pedersen@mail.no"))
			Expect(user.PasswordHash).To(BeEquivalentTo("($)#)@__$()#)@_@($(@)"))
			Expect(user.Teams).To(BeEmpty())
			Expect(user.Permissions).To(BeEmpty())
			users, err := repository.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(users).ToNot(BeEmpty())
		})
	})

	Context("given a user exist", func() {
		identifier := uuid.New()
		AfterEach(func() {
			reset(dataSource)
		})

		BeforeEach(func() {
			addPermission("apps.deploy", dataSource)
			addTeam("spartan", dataSource)
			repository.Create(&users.UserEntry{
				Name:         "Nicole Pedersen",
				NickName:     "NicPed",
				Identifier:   identifier,
				Email:        "nicole.pedersen@mail.no",
				PasswordHash: "($)#)@__$()#)@_@($(@)",
			})
		})
		It("can update the user", func() {
			err := repository.Update(&users.UpdateEntry{
				Email:    "nicole.pedersen@mail.no",
				NickName: "Nonik",
			})
			Expect(err).ToNot(HaveOccurred())
			user, err := repository.FindById(identifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(user.NickName).To(BeEquivalentTo("Nonik"))
		})
		It("can add teams to a user", func() {
			addTeam("argus", dataSource)
			err := repository.AddTeams("nicole.pedersen@mail.no", []string{"argus"})
			Expect(err).ToNot(HaveOccurred())
			user, err := repository.FindById(identifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(user.Teams)).To(BeEquivalentTo(1))
			Expect(user.Teams[0].Name).To(BeEquivalentTo("argus"))
		})
		It("can remove teams from a user", func() {
			addTeam("argus", dataSource)
			addTeam("victus", dataSource)
			repository.AddTeams("nicole.pedersen@mail.no", []string{"argus", "victus"})
			user, err := repository.FindById(identifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(user.Teams)).To(BeEquivalentTo(2))
			err = repository.RemoveTeams("nicole.pedersen@mail.no", []string{"argus"})
			Expect(err).ToNot(HaveOccurred())
			user, _ = repository.FindById(identifier)
			Expect(len(user.Teams)).To(BeEquivalentTo(1))
			Expect(user.Teams[0].Name).To(BeEquivalentTo("victus"))
		})
		It("can add permissions to a user", func() {
			addPermission("envs.set", dataSource)
			err := repository.AddPermissions("nicole.pedersen@mail.no", []string{"envs.set"})
			Expect(err).ToNot(HaveOccurred())
			user, _ := repository.FindById(identifier)
			Expect(len(user.Permissions)).To(BeEquivalentTo(1))
			Expect(user.Permissions[0].Name).To(BeEquivalentTo("envs.set"))
		})
		It("can remove permissions from a user", func() {
			addPermission("envs.set", dataSource)
			addPermission("apps.deploy", dataSource)
			err := repository.AddPermissions("nicole.pedersen@mail.no", []string{"envs.set", "apps.deploy"})
			Expect(err).ToNot(HaveOccurred())
			user, _ := repository.FindById(identifier)
			Expect(len(user.Permissions)).To(BeEquivalentTo(2))
			err = repository.RemovePermissions("nicole.pedersen@mail.no", []string{"apps.deploy"})
			Expect(err).ToNot(HaveOccurred())
			user, _ = repository.FindById(identifier)
			Expect(len(user.Permissions)).To(BeEquivalentTo(1))
			Expect(user.Permissions[0].Name).To(BeEquivalentTo("envs.set"))
		})
		It("cannot create a user with the same email", func() {
			err := repository.Create(&users.UserEntry{
				Name:         "Racoon Mike",
				NickName:     "Rakmi",
				Identifier:   uuid.New(),
				Email:        "nicole.pedersen@mail.no",
				PasswordHash: "($)#)@__$()#)@_@($(@)",
			})
			Expect(err).To(HaveOccurred())
		})
		It("can remove a user", func() {
			err := repository.Remove("nicole.pedersen@mail.no")
			Expect(err).ToNot(HaveOccurred())
			users, err := repository.List()
			Expect(err).ToNot(HaveOccurred())
			Expect(users).To(BeEmpty())
		})
	})

})
