package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestJsonParse_Perform(t *testing.T) {
	t.Parallel()
	input := cltest.RunResultWithValue(`{"high": "11850.00", "last": "11779.99", "timestamp": "1512487535", "bid": "11779.89", "vwap": "11525.17", "volume": "12916.67066094", "low": "11100.00", "ask": "11779.99", "open": 11613.07}`)

	adapter := adapters.JsonParse{Path: []string{"last"}}
	result := adapter.Perform(input, nil)
	val, err := result.Value()
	assert.Equal(t, "11779.99", val)
	assert.Nil(t, err)
	assert.Nil(t, result.GetError())
}

func TestJsonParse_Perform_NonExistingPath(t *testing.T) {
	t.Parallel()
	initialValue := `{"high": "11850.00", "last": "11779.99", "timestamp": "1512487535", "bid": "11779.89", "vwap": "11525.17", "volume": "12916.67066094", "low": "11100.00", "ask": "11779.99", "open": 11613.07}`
	input := cltest.RunResultWithValue(initialValue)

	adapter := adapters.JsonParse{Path: []string{"doesnotexist"}}
	result := adapter.Perform(input, nil)
	val, err := result.Get("value")
	assert.False(t, val.Exists())
	assert.Nil(t, err)
	assert.Nil(t, result.GetError())

	adapter = adapters.JsonParse{Path: []string{"doesnotexist", "noreally"}}
	result = adapter.Perform(input, nil)
	valString, err := result.Value()
	assert.Equal(t, initialValue, valString)
	assert.Nil(t, err)
	assert.NotNil(t, result.GetError())
}
