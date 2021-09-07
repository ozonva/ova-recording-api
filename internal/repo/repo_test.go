package repo_test

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/ozonva/ova-recording-api/internal/repo"
	"github.com/ozonva/ova-recording-api/pkg/recording"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var _ = Describe("Repo", func() {
	var (
		db *sqlx.DB
		testRepo repo.Repo
		testData []recording.Appointment
		additionalEntry recording.Appointment
	)

	BeforeEach(func() {
		db = connectDb()

		prepareFixtures(db.DB, baseMigrationsPath, fixturesPath)

		testRepo = repo.NewRepo(db)

		testData = []recording.Appointment{
			{
				AppointmentID: 1,
				UserID: 1,
				Name: "Some appointment 1",
				Description: "Some description 1",
				StartTime: time.Date(2021,time.September,3, 11, 11, 11, 0, time.UTC),
				EndTime: time.Date(2021,time.September,4, 11, 11, 11, 0, time.UTC),
			},
			{
				AppointmentID: 2,
				UserID: 2,
				Name: "Some appointment 2",
				Description: "Some description 2",
				StartTime: time.Date(2021,time.September,5, 11, 11, 11, 0, time.UTC),
				EndTime: time.Date(2021,time.September,6, 11, 11, 11, 0, time.UTC),
			},
			{
				AppointmentID: 3,
				UserID: 3,
				Name: "Some appointment 3",
				Description: "Some description 3",
				StartTime: time.Date(2021,time.September,7, 11, 11, 11, 0, time.UTC),
				EndTime: time.Date(2021,time.September,8, 11, 11, 11, 0, time.UTC),
			},
		}

		additionalEntry = recording.Appointment{
			AppointmentID: 4,
			UserID: 2,
			Name: "noname1",
			Description: "no desc1",
			StartTime: time.Date(2021,time.January,1, 1, 1, 1, 0, time.UTC),
			EndTime: time.Date(2021,time.January,2, 1, 1, 1, 0, time.UTC),
		}
	})

	AfterEach(func() {
		if db != nil {
			err := db.Close()
			if err != nil {
				logrus.Fatalf("Cannot close connection: %s", err)
			}
		}
	})

	Context("general", func() {

		It("Add entity", func() {
			res, err := testRepo.AddEntities(
				context.Background(),
				[]recording.Appointment{additionalEntry})
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(res[0]).To(gomega.BeEquivalentTo(additionalEntry.AppointmentID))
		})

		It("Update entity", func() {
			entityToUpdate := recording.Appointment{
				AppointmentID: testData[0].AppointmentID,
				UserID: 10,
				Name: "UPDATED",
				StartTime: time.Time{},
				EndTime: time.Time{},
			}

			err := testRepo.UpdateEntity(context.Background(), entityToUpdate)
			gomega.Expect(err).To(gomega.BeNil())

			res, err := testRepo.DescribeEntity(context.Background(), entityToUpdate.AppointmentID)

			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(res.UserID).To(gomega.BeEquivalentTo(entityToUpdate.UserID))
			gomega.Expect(res.Name).To(gomega.BeEquivalentTo(entityToUpdate.Name))
			gomega.Expect(res.Description).To(gomega.BeEquivalentTo(testData[0].Description))
			gomega.Expect(res.StartTime).To(gomega.BeEquivalentTo(testData[0].StartTime))
			gomega.Expect(res.EndTime).To(gomega.BeEquivalentTo(testData[0].EndTime))
		})

		It("Remove entity", func() {
			err := testRepo.RemoveEntity(
				context.Background(),
				4)
			gomega.Expect(err).To(gomega.BeNil())
		})

		It("List entities", func() {
			_, err := testRepo.AddEntities(context.Background(), []recording.Appointment{additionalEntry})
			gomega.Expect(err).To(gomega.BeNil())

			res, err := testRepo.ListEntities(
				context.Background(),
				4, 0)
			gomega.Expect(err).To(gomega.BeNil())

			fullTestData := append(testData, additionalEntry)

			gomega.Expect(res).To(gomega.BeEquivalentTo(fullTestData))
		})

		It("Describe entity", func() {
			res, err := testRepo.DescribeEntity(context.Background(), 1)
			gomega.Expect(err).To(gomega.BeNil())

			gomega.Expect(*res).To(gomega.BeEquivalentTo(testData[0]))
		})
	})
})

func connectDb() *sqlx.DB {
	db, err := sqlx.Connect("txdb", "test_database")
	if err != nil {
		logrus.Fatalf("Cannot connect to database: %s", err)
	}
	return db
}

func prepareFixtures(db *sql.DB, baseMigrationsPath string, fixturesPath string) {
	_, err := db.Exec("CREATE TABLE goose_db_version (id serial primary key, version_id bigint, is_applied boolean, tstamp timestamp default now())")
	if err != nil {
		logrus.Fatalf("Cannot create version table %s", err)
	}

	_, err = db.Exec("INSERT INTO goose_db_version (version_id, is_applied) VALUES (1, true)")
	if err != nil {
		logrus.Fatalf("Cannot insert %s", err)
	}

	goose.SetBaseFS(os.DirFS(baseMigrationsPath))
	err = goose.Up(db, "migrations")
	if err != nil {
		logrus.Fatalf("Cannot up migrations: %s", err)
	}

	goose.SetBaseFS(os.DirFS(fixturesPath))
	err = goose.Up(db, "migrations")
	if err != nil {
		logrus.Fatalf("Cannot up test data: %s", err)
	}
}
