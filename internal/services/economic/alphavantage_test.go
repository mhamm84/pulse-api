package economic

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCpiAlphaParse(t *testing.T) {
	jsonData := `{
		"Name":"Consumer Price Index for all Urban Consumers",
		"Interval":"monthly",
		"Unit":"index 1982-1984=100",
		"Data":[
			{"Date":"2022-05-01","Value":"292.296"},
			{"Date":"2022-04-01","Value":"289.109"},
			{"Date":"2022-03-01","Value":"287.504"}
		]
	}`
	byteData := []byte(jsonData)
	var cpiAlphaResponse AlphaVantageEconomicResponse

	err := json.Unmarshal(byteData, &cpiAlphaResponse)
	assert.NoError(t, err, "error occurred marshalling AlphaVantageEconomicResponse")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "Consumer Price Index for all Urban Consumers", cpiAlphaResponse.Name)
	assert.Equal(t, "monthly", cpiAlphaResponse.Interval)
	assert.Equal(t, "index 1982-1984=100", cpiAlphaResponse.Unit)
	assert.Len(t, cpiAlphaResponse.Data, 3)
}
