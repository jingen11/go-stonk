package main

import (
	"log"
	"os"

	"github.com/jingen11/stonk-tracker/internal/command"
	"github.com/jingen11/stonk-tracker/internal/db"
	stonkapi "github.com/jingen11/stonk-tracker/internal/stonkApi"
	"github.com/jingen11/stonk-tracker/internal/utils"
	"github.com/joho/godotenv"
)

type commands struct {
	Commands map[string]func(*command.Command) error
}

func newCommands() commands {
	c := commands{
		Commands: map[string]func(*command.Command) error{},
	}
	return c
}

func (c *commands) register(name string, f func(*command.Command) error) {
	c.Commands[name] = f
}

func (c *commands) run(name string, command *command.Command) error {
	err := c.Commands[name](command)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	cfg := utils.ProjectConfig{}
	cfg.HistoricalTimeFrame = 100
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

	// polygonKey1 := os.Getenv("POLYGON_IO_KEY_1")
	polygonKey2 := os.Getenv("POLYGON_IO_KEY_2")
	cfg.ApiClient = stonkapi.InitStonkApiClient([]string{polygonKey2})
	cfg.Query = &db.Query{
		PriceColl:  priceColl,
		SymbolColl: symbolColl,
	}

	c := newCommands()

	c.register("refresh", command.HandleRefresh)
	c.register("add", command.HandlerAddNewSymbol)

	comm := os.Args[1]
	args := os.Args[2:]
	err = c.run(comm, &command.Command{
		Cfg:   cfg,
		Input: args,
	})
	if err != nil {
		log.Fatalf("error running command: %v", err)
		os.Exit(1)
	}
}
