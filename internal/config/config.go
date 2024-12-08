package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Secret     string           `yaml:"secret"`
	Salt       string           `yaml:"salt"`
	Cloudflare CloudflareConfig `yaml:"cloudflare"`
	MysqlDSN   string           `yaml:"mysqldsn"`
}

type CloudflareConfig struct {
	AccountID       string `yaml:"account_id"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	BucketName      string `yaml:"bucket_name"`
	BucketDomain    string `yaml:"bucket_domain"`
}

var (
	once sync.Once
	conf Config
)

func Get() Config {
	once.Do(func() {
		filePath := "config.yaml"
		// create an empty config yaml if file not exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			newFile, cErr := os.Create(filePath)
			if cErr != nil {
				panic(cErr)
			}
			defer newFile.Close()
			_ = yaml.NewEncoder(newFile).Encode(&Config{})
			panic("please setup config.yaml")
		}

		yamlFile, err := os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
		if err != nil {
			panic(err)
		}
		if err = yaml.NewDecoder(yamlFile).Decode(&conf); err != nil {
			panic(err)
		}

	})

	return conf
}
