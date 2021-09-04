package repo_test

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/jackc/pgx/stdlib"
	. "github.com/onsi/ginkgo"
	api "github.com/ozonva/ova-recording-api/internal/app/recording"
	"github.com/sirupsen/logrus"
)

const (
	fixturesPath = "/home/evyalyy/ozonva/ova-recording-api/fixtures"
	baseMigrationsPath = "/home/evyalyy/ozonva/ova-recording-api/migrations"
)

var _ = Describe("Repo", func() {
	var (
		db *sql.DB
		fixtures *testfixtures.Loader
	)

	BeforeEach(func() {
		cfg := api.ReadConfig(fixturesPath + "/config/db_config.yml")
		db, err := sql.Open("pgx", cfg.GetConnString())
		if err != nil {
			logrus.Fatalf("Cannot connect to database: %s", err)
		}

		fixtures, err = testfixtures.New(
			testfixtures.Database(db),
			testfixtures.Dialect("postgres"),
			testfixtures.Directory(fixturesPath))

		if err != nil {
			logrus.Fatalf("Cannot create fixtures: %s", err)
		}

		err = fixtures.Load()
		if err != nil {
			logrus.Fatalf("Cannot load fixtures: %s", err)
		}
	})

	AfterEach(func() {
		err := db.Close()
		if err != nil {
			logrus.Fatalf("Cannot close connection: %s", err)
		}
	})

	Context("negative chunk size", func() {
			It("should return all entries", func() {

			})
	})
})
