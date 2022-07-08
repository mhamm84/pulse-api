package economic

import (
	"context"
	"github.com/mhamm84/gofinance-alpha/alpha"
	data2 "github.com/mhamm84/gofinance-alpha/alpha/data"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/shopspring/decimal"
	"time"
)

const (
	cpiServiceTimeout = 10
)

type CpiAlphaService struct {
	Models data.Models
	Client *alpha.AlphaClient
	Logger *jsonlog.Logger
}

// CpiGetAll Gets all the CPI data
// if no data is found, a request is sent to the API to get the data to populate the DB
func (s CpiAlphaService) CpiGetAll(ctx context.Context) (*[]data.Cpi, error) {
	cpiPulseData, err := s.Models.CpiModel.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return cpiPulseData, nil
}

func (s CpiAlphaService) StartCpiDataSyncTask() {
	tr := utils.NewScheduleTaskRunner(5*time.Second, 24*time.Hour, s.Logger)
	tr.Start(func() {
		s.Logger.PrintInfo("AlphaCpi | CpiGetAll | DataSyncTask", map[string]interface{}{
			"service":  "CpiAlphaService",
			"function": "CpiGetAll",
		})
		ctx, cancel := context.WithTimeout(context.Background(), cpiServiceTimeout*time.Second)
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
			s.insertNewData(cpiApiData, cpiData, func(data *data.Cpi) error {
				return s.Models.CpiModel.Insert(ctx, data)
			})
		}
	})
}

func (s CpiAlphaService) insertNewData(apiData *[]data.Cpi, mongoData *[]data.Cpi, insert func(data *data.Cpi) error) error {
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
			//s.Models.CpiModel.InsertOne(&newData)
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

func (s CpiAlphaService) getDataFromApi() (*[]data.Cpi, error) {
	//API
	apiRes, err := s.Client.Cpi(nil)
	if err != nil {
		s.Logger.PrintError(err, nil)
	}
	// Transform
	cpiAppData := make([]data.Cpi, 0, 50)
	for _, d := range apiRes.Data {
		cpiAppData = append(cpiAppData, *s.transform(&d))
	}
	return &cpiAppData, nil
}

func (s CpiAlphaService) transform(apiData *data2.CpiDataValue) *data.Cpi {
	cpiDate, err := time.Parse("2006-01-02", apiData.Date)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"cpiAlphaDate": apiData.Date})
	}
	cpiValue, err := decimal.NewFromString(apiData.Value)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"cpiAlphaValue": apiData.Value})
	}
	return &data.Cpi{
		Date:  cpiDate,
		Value: cpiValue,
	}
}

func (s CpiAlphaService) insertMany(toSave *[]data.Cpi) error {
	ctx, cancel := context.WithTimeout(context.Background(), cpiServiceTimeout*time.Second)
	defer cancel()
	err := s.Models.CpiModel.InsertMany(ctx, toSave)
	if err != nil {
		s.Logger.PrintError(err, nil)
		return err
	}
	return nil

}
