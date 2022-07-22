package economic

import (
	"context"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/internal/jsonlog"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
	"time"
)

type MockEconomicRepository struct {
	mock.Mock
}

func (w *MockEconomicRepository) LatestWithPercentChange(ctx context.Context, table string) (*data.EconomicWithChange, error) {
	args := w.Called(ctx, table)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*data.EconomicWithChange), args.Error(1)
}

func (w *MockEconomicRepository) GetIntervalWithPercentChange(ctx context.Context, table string, years int, paging data.Paging) (*data.EconomicWithChangeResult, error) {
	return nil, nil
}
func (w *MockEconomicRepository) GetAll(ctx context.Context, table string) (*[]data.Economic, error) {
	return nil, nil
}
func (w *MockEconomicRepository) Insert(ctx context.Context, table string, data *data.Economic) error {
	return nil
}
func (w *MockEconomicRepository) InsertMany(ctx context.Context, table string, data *[]data.Economic) error {
	return nil
}

func TestDashboardService_GetDashboardSummary(t *testing.T) {

	ctx := context.Background()
	date, _ := time.Parse("2006-06-02", "2022-01-01")
	dashName := "cpi-dash"
	value := decimal.NewFromFloat(100.00)
	change := decimal.NewFromFloat(10.00)
	slug := "cpi"
	logging := jsonlog.New(os.Stdout, jsonlog.LevelDebug)

	t.Run("createDashSummary", func(t *testing.T) {
		mockRepo := new(MockEconomicRepository)
		mockRepo.On("LatestWithPercentChange", mock.Anything, slug).Return(&data.EconomicWithChange{
			Date:   date,
			Value:  value,
			Change: change,
		}, nil).Once()

		res := createDashSummary(ctx, logging, mockRepo, slug, dashName, nil)
		mockRepo.AssertExpectations(t)

		assert.Equal(t, dashName, res.Name)
		assert.Equal(t, value, res.Value)
		assert.Equal(t, change, res.Change)
		assert.Equal(t, slug, res.Slug)
		assert.Nil(t, res.Extras)
	})
}

func TestDashboardService_GetDashboardSummary_Error(t *testing.T) {

	ctx := context.Background()
	dashName := "cpi-dash"
	slug := "cpi"
	logging := jsonlog.New(os.Stdout, jsonlog.LevelDebug)
	error := errors.New("LatestWithPercentChange error")

	t.Run("createDashSummary", func(t *testing.T) {
		mockRepo := new(MockEconomicRepository)
		mockRepo.On("LatestWithPercentChange", mock.Anything, slug).Return(nil, error).Once()

		res := createDashSummary(ctx, logging, mockRepo, slug, dashName, nil)
		mockRepo.AssertExpectations(t)

		assert.Nil(t, res)
	})
}
