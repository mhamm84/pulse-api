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
	"strconv"
	"time"
)

type AlphaVantageConsumerSentimentService struct {
	Models data.Models
	Client *alpha.Client
	Logger *jsonlog.Logger
}

func (s AlphaVantageConsumerSentimentService) ConsumerSentimentGetAll(ctx context.Context) (*[]economic.ConsumerSentiment, error) {
	csData, err := s.Models.ConsumerSentimentModel.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return csData, nil
}

func (s AlphaVantageConsumerSentimentService) StartDataSyncTask() {
	tr := utils.NewScheduleTaskRunner(5*time.Second, 24*time.Hour, s.Logger)
	tr.Start(func() {
		s.Logger.PrintInfo("AlphaVantageConsumerSentimentService | StartDataSyncTask", map[string]interface{}{
			"service": "AlphaVantageConsumerSentimentService",
		})

		ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
		defer cancel()

		data, err := s.ConsumerSentimentGetAll(ctx)
		if err != nil {
			s.Logger.PrintError(err, nil)
			return
		}

		// Initially Empty, Get data from API and insert
		if len(*data) <= 0 {
			s.Logger.PrintInfo("no data found in DB, getting from API", map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			apiData, err := s.getDataFromApi()
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.Logger.PrintInfo("inserting API data into DB", map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			err = s.insertMany(apiData)
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			return
		}
		// If there is data, check to see if the API has new data
		if len(*data) > 0 {
			// Get transformed API data
			s.Logger.PrintInfo("existing Consumer Sentiment data in DB, checking API for updates", map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			apiData, err := s.getDataFromApi()
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.insertNewData(apiData, data, func(data *economic.ConsumerSentiment) error {
				return s.Models.ConsumerSentimentModel.Insert(ctx, data)
			})
		}
	})
}

func (s AlphaVantageConsumerSentimentService) insertNewData(apiData *[]economic.ConsumerSentiment, mongoData *[]economic.ConsumerSentiment, insert func(data *economic.ConsumerSentiment) error) error {
	if len(*apiData) > len(*mongoData) {
		delta := len(*apiData) - len(*mongoData)
		i := 1
		s.Logger.PrintInfo("new Consumer Sentiment data found form AlphaVantage", map[string]interface{}{
			"delta": delta,
		})
		for delta > 0 {
			newData := (*apiData)[len(*apiData)-i]
			s.Logger.PrintInfo("inserting new Consumer Sentiment data point", map[string]interface{}{
				"consumer_sentiment_date":  newData.Date,
				"consumer_sentiment_value": newData.Value,
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
		s.Logger.PrintInfo("no new Consumer Sentiment data found form AlphaVantage", nil)
		return nil
	}
}

func (s AlphaVantageConsumerSentimentService) getDataFromApi() (*[]economic.ConsumerSentiment, error) {
	//API
	apiRes, err := s.Client.ConsumerSentiment(nil)
	if err != nil {
		s.Logger.PrintError(err, nil)
	}
	// Transform
	appData := make([]economic.ConsumerSentiment, 0, 50)
	for _, d := range apiRes.Data {
		transformedData := s.transform(&d)
		if transformedData != nil {
			appData = append(appData, *transformedData)
		}
	}
	return &appData, nil
}

func (s AlphaVantageConsumerSentimentService) transform(apiData *data2.EconomicValue) *economic.ConsumerSentiment {
	date, err := time.Parse("2006-01-02", apiData.Date)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"AlphaVantageDate": apiData.Date})
	}
	if _, err := strconv.Atoi(apiData.Value); err != nil {
		// skip as this value is not a number
		return nil
	}
	value, err := decimal.NewFromString(apiData.Value)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"AlphaVantageValue": apiData.Value})
	}
	return &economic.ConsumerSentiment{
		Date:  date,
		Value: value,
	}
}

func (s AlphaVantageConsumerSentimentService) insertMany(toSave *[]economic.ConsumerSentiment) error {
	ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
	defer cancel()
	err := s.Models.ConsumerSentimentModel.InsertMany(ctx, toSave)
	if err != nil {
		s.Logger.PrintError(err, nil)
		return err
	}
	return nil

}
