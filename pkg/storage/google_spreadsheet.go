package storage

import (
	"bitcoin-app/pkg/explorer"
	"context"
	"encoding/base64"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSpreadsheetConfig struct {
	Credential    string
	SpreadSheetId string
	SheetId       int
}

type GoogleSpreadsheetService struct {
	Context context.Context
	Service *sheets.Service
	Config  GoogleSpreadsheetConfig
}

// Load new google spreadsheet service
func New(c GoogleSpreadsheetConfig) (*GoogleSpreadsheetService, error) {
	// create api context
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	credBytes, err := base64.StdEncoding.DecodeString(c.Credential)
	if err != nil {
		return nil, err
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	service, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	googleSpreadSheetService := &GoogleSpreadsheetService{
		Context: ctx,
		Service: service,
		Config:  c,
	}
	return googleSpreadSheetService, nil
}

// Append newest block information to google spread sheet
func (svc *GoogleSpreadsheetService) AppendBlockInfo(bitqueryResult explorer.BitQueryResult) error {

	spreadsheetId := svc.Config.SpreadSheetId
	sheetId := svc.Config.SheetId
	service := svc.Service
	ctx := svc.Context

	err := svc.removeAllData()
	if err != nil {
		return err
	}

	// Convert sheet ID to sheet name.
	spreadsheet, err := service.Spreadsheets.Get(spreadsheetId).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || spreadsheet.HTTPStatusCode != 200 {
		return err
	}

	sheetName := ""
	for _, v := range spreadsheet.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetId) {
			sheetName = prop.Title
			break
		}
	}

	// Add title of columns
	row := &sheets.ValueRange{
		Values: [][]interface{}{{"Height", "Difficulty", "Timestamp"}},
	}
	result, err := service.Spreadsheets.Values.Append(spreadsheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("OVERWRITE").Context(ctx).Do()
	if err != nil || result.HTTPStatusCode != 200 {
		return err
	}

	blocks := bitqueryResult.Data.Bitcoin.Blocks
	for i := 0; i < len(blocks); i++ {
		block := blocks[i]
		// Append each block value to the sheet.
		row := &sheets.ValueRange{
			Values: [][]interface{}{{block.Height, block.Difficulty, block.Timestamp.Time}},
		}

		result, err := service.Spreadsheets.Values.Append(spreadsheetId, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("OVERWRITE").Context(ctx).Do()
		if err != nil || result.HTTPStatusCode != 200 {
			return err
		}
	}
	return nil
}

// Since google spreadsheet should store new information when user press the 'update' button,
// All remain data would be removed and new data should be added.
func (svc *GoogleSpreadsheetService) removeAllData() error {
	readRange := "A:C"
	_, err := svc.Service.Spreadsheets.Values.Clear(svc.Config.SpreadSheetId, readRange, &sheets.ClearValuesRequest{}).Context(svc.Context).Do()

	if err != nil {
		return err
	}

	return nil
}
