package utils

import (
	"bufio"

	"github.com/jingen11/stonk-tracker/internal/db"
	stonkapi "github.com/jingen11/stonk-tracker/internal/stonkApi"
)

type ProjectConfig struct {
	ApiClient           *stonkapi.StonkApiClient
	Query               *db.Query
	HistoricalTimeFrame int
	StdScanner          *bufio.Scanner
}
