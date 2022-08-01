//go:build integration

package economic

import (
	"github.com/mhamm84/pulse-api/testing/integration/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetCPIData(t *testing.T) {

	t.Run("TestGetCPIData", func(t *testing.T) {
		cleanup := setupTest(t)
		defer cleanup(t)

		assert.Truef(t, true, "")
	})

	t.Run("TestGetRetailSalesData", func(t *testing.T) {
		cleanup := setupTest(t)

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
