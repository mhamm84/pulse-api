//go:build integration

package economic

import (
	"context"
	"encoding/json"
	"github.com/mhamm84/pulse-api/internal/data"
	"github.com/mhamm84/pulse-api/testing/integration/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

var httpClient = &http.Client{Timeout: time.Second * 30}

func TestGetCPIData(t *testing.T) {

	t.Run("TestGetRetailSalesData", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup(t)

		assert.Truef(t, true, "")
	})

	t.Run("TestGetCPIData", func(t *testing.T) {
		cleanup := setupTest(t)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:9081/v1/economic/cpi", nil)
		assert.NoError(t, err, "")

		response, err := httpClient.Do(req)
		assert.NoError(t, err, "")

		defer func() {
			if tmpErr := response.Body.Close(); tmpErr != nil {
				err = tmpErr
			}
		}()

		decoder := json.NewDecoder(response.Body)

		res := data.EconomicWithChangeResult{}
		err = decoder.Decode(&res)
		assert.NoError(t, err, "")

		assert.Truef(t, true, "")
		cleanup(t)
	})
}

func setupTest(t *testing.T) func(t *testing.T) {
	t.Log("Doing TestGetCPIData Migration.Up...")
	err := config.TestingConfig.Migration.Up()
	if err != nil {
		t.Log(err)
	}
	err = config.TestingConfig.InsertTestData()
	if err != nil {
		t.Log(err)
	}

	require.NoError(t, err)

	return func(t *testing.T) {
		t.Log("teardown")
		err = config.TestingConfig.Migration.Down()
		if err != nil {
			t.Log(err)
		}
		require.NoError(t, err)
	}
}
