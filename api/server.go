package api

import (
	"bitcoin-app/pkg/explorer"
	"bitcoin-app/pkg/storage"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type ServerConfig struct {
	Port string
}

type AppConfig struct {
	explorerConfig explorer.BitQueryConfig
	storageConfig  storage.GoogleSpreadsheetConfig
}

var appConfig AppConfig

func New(svcConfig ServerConfig, expConfig explorer.BitQueryConfig, stgConfig storage.GoogleSpreadsheetConfig) {

	// set instance
	e := echo.New()

	// set middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	appConfig = AppConfig{
		explorerConfig: expConfig,
		storageConfig:  stgConfig,
	}

	// set router
	v1 := e.Group("/api/v1")
	{
		v1.GET("/bitcoin", update)
	}

	// start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", svcConfig.Port)))
}

func update(e echo.Context) error {

	// Get limit information from parameter
	param := e.QueryParam("limit")
	limit, err := strconv.Atoi(param)
	if err != nil {
		return err
	}

	appConfig.explorerConfig.Variables.Limit = limit

	bitqueryResult, err := explorer.UpdateBitcoinInfo(appConfig.explorerConfig)
	if err != nil {
		return err
	}

	googleService, err := storage.New(appConfig.storageConfig)
	if err != nil {
		return err
	}

	printString := fmt.Sprintf("%d number of Bitcoin information from BitQuery is updated to google spreadsheet\n", limit)

	googleService.AppendBlockInfo(bitqueryResult)
	return e.JSON(http.StatusOK, printString)
}
