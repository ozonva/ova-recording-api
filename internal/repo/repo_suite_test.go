package repo_test

import (
	"github.com/DATA-DOG/go-txdb"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/ozonva/ova-recording-api/internal/app/recording"
)

const (
	fixturesPath = "/home/evyalyy/ozonva/ova-recording-api/fixtures"
	baseMigrationsPath = "/home/evyalyy/ozonva/ova-recording-api/"
)

func TestRepo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repo Suite")
}

var _ = BeforeSuite(func() {
	cfg := api.ReadConfig(fixturesPath + "/config/db_config.yml")
	txdb.Register("txdb", "pgx", cfg.GetConnString())
})
