package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/leonardobiffi/envctl/util/file"
	"github.com/leonardobiffi/envctl/util/system/exit"
)

// Config for the project.
type Config struct {
	Profile            string            `json:"profile"`
	Region             string            `json:"region"`
	Environments       map[string]string `json:"environments"`
	DefaultEnvironment string            `json:"defaultEnv"`
}

var fileName = "envctl.json"

func parseConfig(configFile []byte) *Config {
	configuration := &Config{}

	err := json.Unmarshal(configFile, configuration)

	if err != nil {
		exit.Error(err)
	}

	return configuration
}

func parseFromEnvironment() *Config {
	configuration := &Config{
		Profile: os.Getenv("AWS_PROFILE"),
		Region:  os.Getenv("AWS_REGION"),
	}

	return configuration
}

func parseFromFile() *Config {
	configFile, _ := ioutil.ReadFile(fileName)

	return parseConfig(configFile)
}

// CreateConfig adds a configuration file to the project.
func CreateConfig(config *Config) {
	configJSON, _ := json.MarshalIndent(config, "", "  ")

	ioutil.WriteFile(fileName, configJSON, 0644)

	fmt.Println("Setup Complete! " + fileName + " has been added to the project.")
}

// GetConfig reads configuration of the project.
func GetConfig() *Config {
	if !file.Exists(fileName) {
		return parseFromEnvironment()
	}

	return parseFromFile()
}
