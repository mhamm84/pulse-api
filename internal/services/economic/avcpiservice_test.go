package economic

import (
	"fmt"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/data/economic"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/shopspring/decimal"
	"os"
	"testing"
	"time"
)

func TestCpiAlphaService_checkDeltas(t *testing.T) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	underTest := AlphaVantageCpiService{
		Models: data.Models{},
		Client: nil,
		Logger: logger,
	}

	newT1 := time.Now().AddDate(0, -1, 0)
	v1, _ := decimal.NewFromString("12.2")

	newT2 := time.Now().AddDate(0, -1, 0)
	v2, _ := decimal.NewFromString("12.2")

	newT3 := time.Now().AddDate(0, -1, 0)
	v3, _ := decimal.NewFromString("12.2")

	apiData := []economic.Cpi{
		{newT3, v3},
		{newT2, v2},
		{newT1, v1},
	}
	mongoData := []economic.Cpi{
		{newT2, v2},
		{newT1, v1},
	}

	underTest.insertNewData(&apiData, &mongoData, func(data *economic.Cpi) error {
		fmt.Println("inserting data into db")
		return nil
	})
}

func TestCpiAlphaService_checkDeltasSameData(t *testing.T) {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	underTest := AlphaVantageCpiService{
		Models: data.Models{},
		Client: nil,
		Logger: logger,
	}

	newT1 := time.Now().AddDate(0, -1, 0)
	v1, _ := decimal.NewFromString("12.2")

	newT2 := time.Now().AddDate(0, -1, 0)
	v2, _ := decimal.NewFromString("12.2")

	newT3 := time.Now().AddDate(0, -1, 0)
	v3, _ := decimal.NewFromString("12.2")

	apiData := []economic.Cpi{
		{newT3, v3},
		{newT2, v2},
		{newT1, v1},
	}
	mongoData := []economic.Cpi{
		{newT3, v3},
		{newT2, v2},
		{newT1, v1},
	}

	underTest.insertNewData(&apiData, &mongoData, func(data *economic.Cpi) error {
		fmt.Println("inserting data into db")
		return nil
	})
}