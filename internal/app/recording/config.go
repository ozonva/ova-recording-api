package recording

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Host string 	`yaml:"host"`
		Port int		`yaml:"port"`
		Database string	`yaml:"database"`
		User string		`yaml:"user"`
		Password string	`yaml:"password"`
	} `yaml:"database"`
	ChunkSize int `yaml:"chunkSize"`
	Capacity int `yaml:"capacity"`
	SaveIntervalSec int `yaml:"saveIntervalSec"`
}

func (cfg Config) GetConnString() string {
	return fmt.Sprintf("host=%s port=%d database=%s user=%s password='%s'",
						cfg.Database.Host,
						cfg.Database.Port,
						cfg.Database.Database,
						cfg.Database.User,
						cfg.Database.Password)
}

func ReadConfig(configPath string) Config {
	f, err := os.Open(configPath)
	if err != nil {
		logrus.Fatalf("Cannot open config file %s: %s\n", configPath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logrus.Warnf("cannot close file: %s", err)
		}
	}(f)
	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatalf("Cannot parse config file: %s", err)
	}

	return cfg
}
