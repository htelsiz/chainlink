package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestHttpAdapters_NotAUrlError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		adapter adapters.Adapter
	}{
		{"HttpGet", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
		{"HttpPost", &adapters.HttpGet{URL: cltest.MustParseWebURL("NotAURL")}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.adapter.Perform(models.RunResult{}, nil)
			assert.Equal(t, models.JSON{}, result.Data)
			assert.NotNil(t, result.Error)
		})
	}
}

func TestHttpGet_Perform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
	}{
		{"success", 200, "results!", false, `results!`},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`},
		{"not found", 400, "inputValue", true, `<html>so bad</html>`},
		{"server error", 400, "inputValue", true, `Invalid request`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.RunResultWithValue("inputValue")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "GET", test.response,
				func(body string) { assert.Equal(t, ``, body) })
			defer cleanup()

			hga := adapters.HttpGet{URL: cltest.MustParseWebURL(mock.URL)}
			result := hga.Perform(input, store)

			val, err := result.Value()
			assert.Nil(t, err)
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}

func TestHttpPost_Perform(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		status      int
		want        string
		wantErrored bool
		response    string
	}{
		{"success", 200, "results!", false, `results!`},
		{"success but error in body", 200, `{"error": "results!"}`, false, `{"error": "results!"}`},
		{"success with HTML", 200, `<html>results!</html>`, false, `<html>results!</html>`},
		{"not found", 400, "inputVal", true, `<html>so bad</html>`},
		{"server error", 500, "inputVal", true, `big error`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.RunResultWithValue("inputVal")
			mock, cleanup := cltest.NewHTTPMockServer(t, test.status, "POST", test.response,
				func(body string) { assert.Equal(t, `{"value":"inputVal"}`, body) })
			defer cleanup()

			hpa := adapters.HttpPost{URL: cltest.MustParseWebURL(mock.URL)}
			result := hpa.Perform(input, store)

			val, err := result.Get("value")
			assert.Nil(t, err)
			assert.Equal(t, test.want, val.String())
			assert.Equal(t, true, val.Exists())
			assert.Equal(t, test.wantErrored, result.HasError())
			assert.Equal(t, false, result.Pending)
		})
	}
}
