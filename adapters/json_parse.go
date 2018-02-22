package adapters

import (
	"errors"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
)

// JsonParse holds a path to the desired field in a JSON object,
// made up of an array of strings.
type JsonParse struct {
	Path []string `json:"path"`
}

// Perform returns the value associated to the desired field for a
// given JSON object.
//
// For example, if the JSON data looks like this:
//   {
//     "data": [
//       {"last": "1111"},
//       {"last": "2222"}
//     ]
//   }
//
// Then ["0","last"] would be the path, and "111" would be the returned value
func (jpa *JsonParse) Perform(input models.RunResult, _ *store.Store) models.RunResult {
	val, err := input.Value()
	if err != nil {
		return input.WithError(err)
	}

	js, err := simplejson.NewJson([]byte(val))
	if err != nil {
		return input.WithError(err)
	}

	js, err = checkEarlyPath(js, jpa.Path)
	if err != nil {
		return input.WithError(err)
	}

	rval, ok := js.CheckGet(jpa.Path[len(jpa.Path)-1])
	if !ok {
		return models.RunResult{}
	}

	return input.WithValue(rval.MustString())
}

func checkEarlyPath(js *simplejson.Json, path []string) (*simplejson.Json, error) {
	var ok bool
	for _, k := range path[:len(path)-1] {
		js, ok = js.CheckGet(k)
		if !ok {
			return js, errors.New("No value could be found for the key '" + k + "'")
		}
	}
	return js, nil
}
