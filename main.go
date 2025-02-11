package main

import (
	"log"
	"os"

	"github.com/jingen11/stonk-tracker/internal/db"
	stonkapi "github.com/jingen11/stonk-tracker/internal/stonkApi"
	"github.com/jingen11/stonk-tracker/internal/utils"
	"github.com/joho/godotenv"
)

func main() {
	cfg := utils.ProjectConfig{}
	godotenv.Load()
	dbUrl := os.Getenv("MONGODB_URL")

	dbClient, err := db.Init(dbUrl)
	defer db.Disconnect(dbClient)

	if err != nil {
		log.Fatalf("failed to start mongodb server, error: %s", err.Error())
		os.Exit(1)
	}
	stonkDb := dbClient.Database("stonk")
	priceColl, err := db.InitPriceCollection(stonkDb)
	if err != nil {
		log.Fatalf("failed to intialise price collection, error: %s", err.Error())
		os.Exit(1)
	}
	symbolColl, err := db.InitSymbolCollection(stonkDb)
	if err != nil {
		log.Fatalf("failed to intialise symbol collection, error: %s", err.Error())
		os.Exit(1)
	}

	alphaKey1 := os.Getenv("ALPHA_VANTAGE_KEY_1")
	alphaKey2 := os.Getenv("ALPHA_VANTAGE_KEY_2")
	cfg.ApiClient = stonkapi.InitStonkApiClient([]string{alphaKey1, alphaKey2})
	cfg.Query = &db.Query{
		PriceColl:  priceColl,
		SymbolColl: symbolColl,
	}
}
