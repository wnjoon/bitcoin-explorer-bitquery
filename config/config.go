package config

import (
	"bitcoin-app/api"
	"bitcoin-app/pkg/explorer"
	"bitcoin-app/pkg/storage"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	viper viper.Viper
}

var AppConfig *Config

func New(path, fileName, fileType string) error {
	viper.AddConfigPath(path)     // "."
	viper.SetConfigName(fileName) // "config"
	viper.SetConfigType(fileType) // "yaml"

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	AppConfig = &Config{
		viper: *viper.GetViper(),
	}
	return nil
}

func (c Config) ServerConfig() api.ServerConfig {
	path := "server"
	config := api.ServerConfig{}
	config.Port = c.viper.GetString(path + ".port")
	return config
}

func (c Config) ExplorerConfig() explorer.BitQueryConfig {
	conf := c.bitQueryExplorer("explorer")
	return conf
}

func (c Config) bitQueryExplorer(basePath string) explorer.BitQueryConfig {
	path := basePath + ".bitquery"

	config := explorer.BitQueryConfig{}
	config.URL = c.viper.GetString(path + ".url")
	config.API_KEY = c.viper.GetString(path + ".apikey")

	variables := explorer.Variables{}
	variables.Offset = c.viper.GetInt(path + ".variables.offset")
	variables.Network = c.viper.GetString(path + ".variables.network")
	variables.DateFormat = c.viper.GetString(path + ".variables.dateformat")

	config.Variables = variables
	return config
}

func (c Config) StorageConfig() storage.GoogleSpreadsheetConfig {
	conf := c.googleSpreadsheet("storage")
	return conf
}

func (c Config) googleSpreadsheet(basePath string) storage.GoogleSpreadsheetConfig {
	path := basePath + ".google.spreadsheet"
	config := storage.GoogleSpreadsheetConfig{}
	config.Credential = c.viper.GetString(path + ".credential")
	config.SpreadSheetId = c.viper.GetString(path + ".spreadsheetid")
	config.SheetId = c.viper.GetInt(path + ".sheetid")
	return config
}
