package config

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	Srv     *sheets.Service
	SheetID string
)

func InitSheets() {
	SheetID = os.Getenv("SHEET_ID")
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		panic(err)
	}

	config, err := google.JWTConfigFromJSON(b, sheets.SpreadsheetsScope)
	if err != nil {
		panic(err)
	}

	client := config.Client(ctx)
	Srv, err = sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		panic(err)
	}
}