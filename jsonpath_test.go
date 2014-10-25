package jsonpath

import (
  "encoding/json"
  "strings"
  "testing"

  "github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
  tt := []struct {
    jsn    string
    path   string
    expVal interface{}
  }{
    {`{"a": 1}`, "a", 1},
    {`{"a": {"b": {"c": 2}}}`, "a.b.c", 2},
    {`{"a": [{}, {"b": "wow"}] }`, "a[1].b", "wow"},
    {`{"a": [{},{},{},{},{},{},{},{},{},{},{"b": 1}]}`, "a[10].b", 1},
    {`{"a": [1]}`, "a", []interface{}{1}},
    {`{"a": [1]}`, "a[0]", 1},
  }

  for _, row := range tt {
    var jsn interface{}
    err := json.Unmarshal([]byte(row.jsn), &jsn)
    assert.Nil(t, err)
    got, err := Get(jsn, row.path)
    assert.Nil(t, err)
    assert.Equal(t, got, row.expVal, "got: %s, expected: %s", got, row.expVal)
  }
}

func TestGetErrors(t *testing.T) {
  tt := []struct {
    jsn    string
    path   string
    expErr string
  }{
    {`{"a": 1}`, "a.b", "error finding path"},
    {`{"a": 1}`, "a.", "path empty"},
    {`{"a": 1}`, "a[10]", "did not find expected slice"},
    {`{"a": 1}`, "a[a]", "parsing"},
    {`{"a": {"b": {"c": 2}}}`, "a.c", "error finding path"},
    {`{"a": [{}, {"b": 2}] }`, "a[10].b", "error index out of bound"},
    {`{"a": [{}, {"b": 2}] }`, "a[1b", "error closing bracket missing"},
  }

  for _, row := range tt {
    var jsn interface{}
    err := json.Unmarshal([]byte(row.jsn), &jsn)
    assert.Nil(t, err)
    _, err = Get(jsn, row.path)
    assert.NotNil(t, err)
    assert.True(t, strings.Contains(err.Error(), row.expErr),
      "got: %s, expected: %s", err.Error(), row.expErr)
  }
}
