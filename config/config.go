package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func LoadConfig[T any](target *T, filename string) {
	configData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("error while reading yaml file, %s\n", err)
	}

	if err = yaml.Unmarshal(configData, &target); err != nil {
		log.Fatalf("error while parsing yaml file, %s\n", err)
	}
}
