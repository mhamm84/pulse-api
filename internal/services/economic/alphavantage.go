package economic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhamm84/gofinance-alpha/alpha"
	data2 "github.com/mhamm84/gofinance-alpha/alpha/data"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/data/economic"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	serviceTimeout                   = 10
	cpiTableName                     = "cpi"
	consumerSentimentTableName       = "consumer_sentiment"
	treasuryYieldThreeMonthTableName = "treasury_yield_three_month"
	treasuryYieldTwoYearTableName    = "treasury_yield_two_year"
	treasuryYieldFiveYearTableName   = "treasury_yield_five_year"
	treasuryYieldSevenYearTableName  = "treasury_yield_seven_year"
	treasuryYieldTenYearTableName    = "treasury_yield_ten_year"
	treasuryYieldThirtyYearTableName = "treasury_yield_thirty_year"
)

type alphaEconomicCall func(opts *alpha.Options) (*data2.EconomicResponse, error)
type getAllEconomic func(ctx context.Context, tableName string) (*[]economic.Economic, error)

type AlphaVantageEconomicResponse struct {
	Name     string                     `json:"Name"`
	Interval string                     `json:"Interval"`
	Unit     string                     `json:"Unit"`
	Data     []AlphaVantageEconomicData `json:"Data"`
}

type AlphaVantageEconomicData struct {
	Date  time.Time `json:"date"`
	Value big.Float `json:"value"`
}

func (l *AlphaVantageEconomicData) UnmarshalJSON(j []byte) error {
	var rawStrings map[string]string

	err := json.Unmarshal(j, &rawStrings)
	if err != nil {
		return err
	}

	for k, v := range rawStrings {
		if strings.ToLower(k) == "date" {
			t, err := time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
			l.Date = t
		}
		if strings.ToLower(k) == "value" {
			fv, err := strconv.ParseFloat(v, 64)
			v := big.NewFloat(fv)
			if err != nil {
				return err
			}
			l.Value = *v
		}
	}
	return nil
}

type AlphaVantageEconomicService struct {
	Models data.Models
	Client *alpha.Client
	Logger *jsonlog.Logger
}

// GetAll Gets all the data for an economic table
// if no data is found, a request is sent to the API to get the data to populate the DB
func (s AlphaVantageEconomicService) GetAll(ctx context.Context, tableName string) (*[]economic.Economic, error) {
	data, err := s.Models.EconomicModel.GetAll(ctx, tableName)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s AlphaVantageEconomicService) StartDataSyncTask() {
	s.start(nil, s.Client.Cpi, cpiTableName, s.Models.EconomicModel.GetAll)
	s.start(nil, s.Client.ConsumerSentiment, consumerSentimentTableName, s.Models.EconomicModel.GetAll)
	s.addTreasuryYields()
}

func (s AlphaVantageEconomicService) addTreasuryYields() {
	threeMonthOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.ThreeMonth}
	s.start(&threeMonthOptions, s.Client.TreasuryYield, treasuryYieldThreeMonthTableName, s.Models.EconomicModel.GetAll)
	twoYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.TwoYear}
	s.start(&twoYearOptions, s.Client.TreasuryYield, treasuryYieldTwoYearTableName, s.Models.EconomicModel.GetAll)
	fiveYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.FiveYear}
	s.start(&fiveYearOptions, s.Client.TreasuryYield, treasuryYieldFiveYearTableName, s.Models.EconomicModel.GetAll)
	sevenYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.SevenYear}
	s.start(&sevenYearOptions, s.Client.TreasuryYield, treasuryYieldSevenYearTableName, s.Models.EconomicModel.GetAll)
	tenYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.TenYear}
	s.start(&tenYearOptions, s.Client.TreasuryYield, treasuryYieldTenYearTableName, s.Models.EconomicModel.GetAll)
	thirtyYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.ThirtyYear}
	s.start(&thirtyYearOptions, s.Client.TreasuryYield, treasuryYieldThirtyYearTableName, s.Models.EconomicModel.GetAll)
}

func (s AlphaVantageEconomicService) start(opts *alpha.Options, apiCall alphaEconomicCall, tableName string, getAll getAllEconomic) {
	tr := utils.NewScheduleTaskRunner(5*time.Second, 24*time.Hour, s.Logger)
	tr.Start(func() {
		s.Logger.PrintInfo("StartDataSyncTask", map[string]interface{}{
			"data": tableName,
		})

		ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
		defer cancel()

		data, err := getAll(ctx, tableName)
		if err != nil {
			s.Logger.PrintError(err, nil)
			return
		}

		// Initially Empty, Get data from API and insert
		if len(*data) <= 0 {
			s.Logger.PrintInfo(fmt.Sprintf("no data found in DB for %s, getting from API", tableName), map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			apiData, err := s.getDataFromApi(opts, apiCall)
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.Logger.PrintInfo(fmt.Sprintf("inserting %s API data into DB", tableName), map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			err = s.insertMany(apiData, tableName)
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			return
		}
		// If there is data, check to see if the API has new data
		if len(*data) > 0 {
			// Get transformed API data
			s.Logger.PrintInfo(fmt.Sprintf("existing %s data in DB, checking API for updates", tableName), map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			apiData, err := s.getDataFromApi(opts, apiCall)
			if err != nil {
				s.Logger.PrintError(err, nil)
				return
			}
			s.insertNewData(ctx, tableName, apiData, data)
		}
	})
}

func (s AlphaVantageEconomicService) insertNewData(
	ctx context.Context,
	tableName string,
	apiData *[]economic.Economic,
	dbData *[]economic.Economic) error {

	dbMap := make(map[int64]*economic.Economic)
	for _, data := range *dbData {
		dbMap[data.Date.Unix()] = &data
	}

	for _, data := range *apiData {
		if check := dbMap[data.Date.Unix()]; check == nil {
			s.Logger.PrintInfo(fmt.Sprintf("inserting new data point for %s", tableName), map[string]interface{}{
				"date":  data.Date,
				"value": data.Value,
			})
			err := s.Models.EconomicModel.Insert(ctx, tableName, &data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s AlphaVantageEconomicService) getDataFromApi(opts *alpha.Options, apiCall alphaEconomicCall) (*[]economic.Economic, error) {
	apiRes, err := apiCall(opts)
	if err != nil {
		s.Logger.PrintError(err, nil)
	}
	// Transform
	data := make([]economic.Economic, 0, 50)
	for _, d := range apiRes.Data {
		transformedData := s.transform(&d)
		if transformedData != nil {
			data = append(data, *transformedData)
		}
	}
	return &data, nil
}
func (s AlphaVantageEconomicService) transform(apiData *data2.EconomicValue) *economic.Economic {
	date, err := time.Parse("2006-01-02", apiData.Date)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"AlphaVantage Date": apiData.Date})
	}
	if _, err := strconv.ParseFloat(apiData.Value, 64); err != nil {
		s.Logger.PrintDebug("cannot parse this value: "+apiData.Value, nil)
		return nil
	}
	value, err := decimal.NewFromString(strings.TrimSpace(apiData.Value))
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"AlphaVantage Value": apiData.Value})
	}
	return &economic.Economic{
		Date:  date,
		Value: value,
	}
}

func (s AlphaVantageEconomicService) insertMany(toSave *[]economic.Economic, tableName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
	defer cancel()
	err := s.Models.EconomicModel.InsertMany(ctx, tableName, toSave)
	if err != nil {
		s.Logger.PrintError(err, nil)
		return err
	}
	return nil

}
