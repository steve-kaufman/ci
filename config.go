package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessKey  string      `yaml:"accessKey"`
	Containers []Container `yaml:"containers"`
}

type Container types.ContainerCreateConfig

func (c Config) GetContainer(imageName string) (Container, bool) {
	for _, container := range c.Containers {
		if imageName == container.Config.Image {
			return container, true
		}
	}
	return Container{}, false
}

func loadConfig() Config {
	configFilePath := os.Getenv("CI_CONFIG_FILE")
	configYAML, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error reading YAML config file")
		panic(err)
	}

	var config Config
	yaml.Unmarshal(configYAML, &config)

	return config
}
