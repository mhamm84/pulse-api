package alpha

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mhamm84/gofinance-alpha/alpha"
	alphavantage "github.com/mhamm84/gofinance-alpha/alpha/data"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/mhamm84/pulse-api/internal/utils"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"golang.org/x/time/rate"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const serviceTimeout = 10
const dataSyncTimeout = 30
const taskDelay = 24 * time.Hour

type ClientInterface interface {
	EconomicData(ctx context.Context, reportType alpha.ReportType, opts *alpha.Options) (*alphavantage.EconomicResponse, error)
}

type alphaEconomicCall func(ctx context.Context, reportType alpha.ReportType, opts *alpha.Options) (*alphavantage.EconomicResponse, error)
type economicDataCall func(ctx context.Context, tableName string) (*[]data.Economic, error)

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
	EconomicRepository data.EconomicRepository
	ReportRepository   data.ReportRepository
	Client             ClientInterface
	Logger             *jsonlog.Logger
	Limiter            AlphaVantageLimiter
}

type AlphaVantageLimiter struct {
	MinuteLimiter *rate.Limiter
	DailyLimiter  *rate.Limiter
}

func (s AlphaVantageEconomicService) GetIntervalWithPercentChange(ctx context.Context, dataChan chan data.EconomicWithChangeResult, errChan chan error, reportType data.ReportType, years int, paging data.Paging) {
	data, err := s.EconomicRepository.GetIntervalWithPercentChange(ctx, data.TableFromReportType(reportType), years, paging)
	if err != nil {
		errChan <- err
	} else {
		dataChan <- *data
	}
}

// GetAll Gets all the data for an economic table
// if no data is found, a request is sent to the API to get the data to populate the DB
func (s AlphaVantageEconomicService) GetAll(reportType data.ReportType) (*[]data.Economic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
	defer cancel()

	data, err := s.EconomicRepository.GetAll(ctx, data.TableFromReportType(reportType))
	if err != nil {
		return nil, err
	}
	return data, nil
}

type DataSyncTaskParams struct {
	s          *AlphaVantageEconomicService
	reportType alpha.ReportType
	opts       *alpha.Options
	apiCall    alphaEconomicCall
	tableName  string
	dataCall   economicDataCall
	reportMap  *map[string]data.Report
}

func (s AlphaVantageEconomicService) StartDataSyncTask() {

	ctx, cancel := context.WithTimeout(context.Background(), dataSyncTimeout*time.Second)
	defer cancel()

	reports, err := s.ReportRepository.GetReports(ctx)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{
			"message":  "Could not get economic report info data",
			"service":  "AlphaVantageEconomicService",
			"function": "StartDataSyncTask",
		})
		return
	}
	if reports == nil || len(*reports) == 0 {
		s.Logger.PrintError(err, map[string]interface{}{
			"message":  "No data found for economic report info data",
			"service":  "AlphaVantageEconomicService",
			"function": "StartDataSyncTask",
		})
		return
	}

	reportMap := map[string]data.Report{}
	for _, v := range *reports {
		reportMap[v.Slug] = v
	}

	start(DataSyncTaskParams{&s, alpha.CPI, nil, s.Client.EconomicData, data.TableFromReportType(data.CPI), s.EconomicRepository.GetAll, &reportMap})
	start(DataSyncTaskParams{&s, alpha.CONSUMER_SENTIMENT, nil, s.Client.EconomicData, data.TableFromReportType(data.ConsumerSentiment), s.EconomicRepository.GetAll, &reportMap})
	addTreasuryYields(&s, &reportMap)
	start(DataSyncTaskParams{&s, alpha.RETAIL_SALES, nil, s.Client.EconomicData, data.TableFromReportType(data.RetailSales), s.EconomicRepository.GetAll, &reportMap})
}

func addTreasuryYields(s *AlphaVantageEconomicService, reportMap *map[string]data.Report) {
	threeMonthOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.ThreeMonth}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &threeMonthOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldThreeMonth), s.EconomicRepository.GetAll, reportMap})
	twoYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.TwoYear}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &twoYearOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldTwoYear), s.EconomicRepository.GetAll, reportMap})
	fiveYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.FiveYear}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &fiveYearOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldFiveYear), s.EconomicRepository.GetAll, reportMap})
	sevenYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.SevenYear}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &sevenYearOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldSevenYear), s.EconomicRepository.GetAll, reportMap})
	tenYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.TenYear}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &tenYearOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldTenYear), s.EconomicRepository.GetAll, reportMap})
	thirtyYearOptions := alpha.Options{Interval: alpha.Daily, Maturity: alpha.ThirtyYear}
	start(DataSyncTaskParams{s, alpha.TREASURY_YIELD, &thirtyYearOptions, s.Client.EconomicData, data.TableFromReportType(data.TreasuryYieldThirtyYear), s.EconomicRepository.GetAll, reportMap})
}

func start(taskParams DataSyncTaskParams) {

	s := taskParams.s
	tableName := taskParams.tableName
	opts := taskParams.opts
	apiCall := taskParams.apiCall
	reportMap := *taskParams.reportMap

	report := reportMap[tableName]
	reportDelay := report.InitialSyncDelayMinutes

	var initialDelayDuration time.Duration
	if reportDelay == 0 {
		initialDelayDuration = 5 * time.Second
	} else {
		initialDelayDuration = time.Duration(reportDelay) * time.Minute
	}
	tr := utils.NewScheduleTaskRunner(initialDelayDuration, 24*time.Hour, s.Logger)
	taskParams.s.Logger.PrintInfo("created new ScheduleTaskRunner", map[string]interface{}{
		"report":               tableName,
		"initialDelayDuration": initialDelayDuration.String(),
		"taskDelay":            taskDelay.String(),
	})
	tr.Start(func() {

		if time.Since(reportMap[tableName].LastPullDate) < (time.Hour * 24) {
			s.Logger.PrintInfo("skipping data sync, less than 24 hours since last one", map[string]interface{}{
				"data": tableName,
			})
			return
		}

		s.Logger.PrintInfo("StartDataSyncTask", map[string]interface{}{
			"data": tableName,
		})

		ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
		defer cancel()

		data, err := taskParams.dataCall(ctx, taskParams.tableName)
		if err != nil {
			s.Logger.PrintError(err, nil)
			return
		}
		// Initially Empty, Get data from API and insert
		if len(*data) <= 0 {
			s.Logger.PrintInfo(fmt.Sprintf("no data found in DB for %s, getting from API", tableName), map[string]interface{}{
				"task": "StartDataSyncTask",
			})
			apiData := processApiCall(ctx, s, taskParams.reportType, opts, apiCall, tableName)
			if apiData == nil {
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
			apiData := processApiCall(ctx, s, taskParams.reportType, opts, apiCall, tableName)

			if apiData != nil {
				s.insertNewData(ctx, tableName, apiData, data)
			}
		}
	})
}

func processApiCall(ctx context.Context, s *AlphaVantageEconomicService, reportType alpha.ReportType, opts *alpha.Options, apiCall alphaEconomicCall, tableName string) *[]data.Economic {
	apiData, err := s.getDataFromApi(ctx, reportType, opts, apiCall)
	if err != nil {
		s.Logger.PrintWarning("error getting data from Alpha Vantage API", map[string]interface{}{
			"tableName": tableName,
			"error":     err.Error(),
		})
		return nil
	}
	if apiData == nil || len(*apiData) == 0 {
		s.Logger.PrintWarning("alpha Vantage API call returned empty", map[string]interface{}{
			"data to extract": tableName,
		})
		return nil
	}
	err = s.ReportRepository.UpdateReportLastPullDate(ctx, tableName)
	if err != nil {
		s.Logger.PrintWarning("error updating last data pull date on report", map[string]interface{}{
			"report": tableName,
			"error":  err.Error(),
		})
	}
	return apiData
}

func (s AlphaVantageEconomicService) insertNewData(ctx context.Context, tableName string, apiData *[]data.Economic, dbData *[]data.Economic) error {

	dbMap := make(map[int64]*data.Economic)
	for _, data := range *dbData {
		dbMap[data.Date.Unix()] = &data
	}

	for _, data := range *apiData {
		if check := dbMap[data.Date.Unix()]; check == nil {
			s.Logger.PrintInfo(fmt.Sprintf("inserting new data point for %s", tableName), map[string]interface{}{
				"date":  data.Date,
				"value": data.Value,
			})
			err := s.EconomicRepository.Insert(ctx, tableName, &data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s AlphaVantageEconomicService) getDataFromApi(ctx context.Context, reportType alpha.ReportType, opts *alpha.Options, apiCall alphaEconomicCall) (*[]data.Economic, error) {
	// Check the API limits
	if !s.Limiter.DailyLimiter.Allow() {
		return nil, errors.New("hit daily limit when calling Alpha Advantage")
	}
	if !s.Limiter.MinuteLimiter.Allow() {
		return nil, errors.New("hit minute limit when calling Alpha Advantage")
	}

	apiRes, err := apiCall(ctx, reportType, opts)
	if err != nil {
		s.Logger.PrintError(err, nil)
	}
	// Transform
	data := make([]data.Economic, 0, 50)
	for _, d := range apiRes.Data {
		transformedData := s.transform(&d)
		if transformedData != nil {
			data = append(data, *transformedData)
		}
	}
	return &data, nil
}
func (s AlphaVantageEconomicService) transform(apiData *alphavantage.EconomicValue) *data.Economic {
	date, err := time.Parse("2006-01-02", apiData.Date)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"alpha vantage Date": apiData.Date})
	}
	if _, err := strconv.ParseFloat(apiData.Value, 64); err != nil {
		s.Logger.PrintDebug("cannot parse this value: "+apiData.Value, nil)
		return nil
	}
	value, err := decimal.NewFromString(strings.TrimSpace(apiData.Value))
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{"alpha vantage Value": apiData.Value})
	}
	return &data.Economic{
		Date:  date,
		Value: value,
	}
}

func (s AlphaVantageEconomicService) insertMany(toSave *[]data.Economic, tableName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), serviceTimeout*time.Second)
	defer cancel()
	err := s.EconomicRepository.InsertMany(ctx, tableName, toSave)
	if err != nil {
		s.Logger.PrintError(err, map[string]interface{}{
			"tableName": tableName,
		})
		return err
	}
	return nil
}
