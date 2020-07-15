package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port		 string `json:"port"`
	Version	 	 string `json:"version"`
	Database struct {
		Host     string `json:"host"`
		User	 string `json:"user"`
		Password string `json:"password"`
		Port 	 string `json:"port"`
		Database string `json:"database"`
	} `json:"database"`
}

type ConfigService struct {
	Config *Config
}

const fileName = "./config/config.json"

func CreateConfigService() *ConfigService{
	Config := loadConfig()
	return NewConfigService(Config)
}

func NewConfigService(Config *Config) *ConfigService {
	return &ConfigService{Config}
}

func loadConfig() *Config {
	var config Config

	configFile, err := os.Open(fileName)
	defer configFile.Close()

	if err != nil {
		log.Fatal(err)
	}

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	return &config
}

/*
func (c *Config) ChangeConfig(newVal string) {
	(*c).Database.Database = "test1"
}
*/