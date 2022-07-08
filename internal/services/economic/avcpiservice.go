package economic

import (
	"context"
	"github.com/mhamm84/gofinance-alpha/alpha"
	data2 "github.com/mhamm84/gofinance-alpha/alpha/data"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/data/economic"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/shopspring/decimal"
	"time"
)

type AlphaVantageCpiService struct {
	Models data.Models
	Client *alpha.Client
	Logger *jsonlog.Logger
}

// CpiGetAll Gets all the CPI data
// if no data is found, a request is sent to the API to get the data to populate the DB
func (s AlphaVantageCpiService) CpiGetAll(ctx context.Context) (*[]economic.Cpi, error) {
	cpiPulseData, err := s.Models.CpiModel.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return cpiPulseData, nil
}

func (s AlphaVantageCpiService) StartDataSyncTask() {
	tr := utils.NewScheduleTaskRunner(5*time.Second, 24*time.Hour, s.Logger)
	tr.Start(func() {
		s.Logger.PrintInfo("AlphaVantageCpiService | StartDataSyncTask", map[string]interface{}{
			"service": "AlphaVantageCpiService",
		})

		ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
		defer cancel()

		cpiData, err := s.CpiGetAll(ctx)
		if err != nil {
			s.Logger.PrintError(err, nil)
			return
		}

		// Initially Empty, Get data from API and insert
		if len(*cpiData) <= 0 {
			s.Logger.PrintInfo("no data found in DB, getting from API", map[string]interface{}{
				"task": "StartCpiDataSyncTask",
			})
			cpiApiData, err := s.getDataFromApi()
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.Logger.PrintInfo("inserting API data into DB", map[string]interface{}{
				"task": "StartCpiDataSyncTask",
			})
			err = s.insertMany(cpiApiData)
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			return
		}
		// If there is data, check to see if the API has new data
		if len(*cpiData) > 0 {
			// Get transformed API data
			s.Logger.PrintInfo("existing CPI data in DB, checking API for updates", map[string]interface{}{
				"task": "StartCpiDataSyncTask",
			})
			cpiApiData, err := s.getDataFromApi()
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.insertNewData(cpiApiData, cpiData, func(data *economic.Cpi) error {
				return s.Models.CpiModel.Insert(ctx, data)
			})
		}
	})
}

func (s AlphaVantageCpiService) insertNewData(apiData *[]economic.Cpi, mongoData *[]economic.Cpi, insert func(data *economic.Cpi) error) error {
	if len(*apiData) > len(*mongoData) {
		delta := len(*apiData) - len(*mongoData)
		i := 1
		s.Logger.PrintInfo("new CPI data found form AlphaVantage", map[string]interface{}{
			"delta": delta,
		})
		for delta > 0 {
			newData := (*apiData)[len(*apiData)-i]
			s.Logger.PrintInfo("inserting new CPI data point", map[string]interface{}{
				"cpi_date":  newData.Date,
				"cpi_value": newData.Value,
			})
			err := insert(&newData)
			if err != nil {
				return err
			}
			delta--
			i++
		}
		return nil
	} else {
		s.Logger.PrintInfo("no new CPI data found form AlphaVantage", nil)
		return nil
	}
}

func (s AlphaVantageCpiService) getDataFromApi() (*[]economic.Cpi, error) {
	//API
	apiRes, err := s.Client.Cpi(nil)
	if err != nil {
		s.Logger.PrintError(err, nil)
	}
	// Transform
	cpiAppData := make([]economic.Cpi, 0, 50)
	for _, d := range apiRes.Data {
		cpiAppData = append(cpiAppData, *s.transform(&d))
	}
	return &cpiAppData, nil
}

func (s AlphaVantageCpiService) transform(apiData *data2.EconomicValue) *economic.Cpi {
	cpiDate, err := time.Parse("2006-01-02", apiData.Date)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"cpiAlphaDate": apiData.Date})
	}
	cpiValue, err := decimal.NewFromString(apiData.Value)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"cpiAlphaValue": apiData.Value})
	}
	return &economic.Cpi{
		Date:  cpiDate,
		Value: cpiValue,
	}
}

func (s AlphaVantageCpiService) insertMany(toSave *[]economic.Cpi) error {
	ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
	defer cancel()
	err := s.Models.CpiModel.InsertMany(ctx, toSave)
	if err != nil {
		s.Logger.PrintError(err, nil)
		return err
	}
	return nil

}
