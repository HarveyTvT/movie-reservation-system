package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Secret string
}

var (
	once sync.Once
	conf Config
)

func Get() Config {
	once.Do(func() {
		filePath := "/var/local/movie-reservation/config.yaml"
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.MkdirAll("/var/local/movie-reservation", os.ModePerm)
			os.Create(filePath)
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
