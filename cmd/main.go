package main

import (
	"bitcoin-app/api"
	"bitcoin-app/config"
	"flag"
	"log"
)

func main() {
	var configFileName, configFilePath string

	flag.StringVar(&configFileName, "conf-file", "config", "Config File Name")
	flag.StringVar(&configFilePath, "conf-path", "./config", "Config File Path")
	flag.Parse()

	// load configs from file
	err := config.New(configFilePath, configFileName, "yaml")
	if err != nil {
		log.Fatalln("failed to load config file - ", err.Error())
	}

	// Set configuration
	explorerConfig := config.AppConfig.ExplorerConfig()
	storageConfig := config.AppConfig.StorageConfig()
	serverConfig := config.AppConfig.ServerConfig()

	// Start server
	api.New(serverConfig, explorerConfig, storageConfig)
}
