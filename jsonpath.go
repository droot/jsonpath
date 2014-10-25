// Package jsonpath provides a simple interface to access fields in a JSON blob.
package jsonpath

import (
  "fmt"
  "strconv"
  "strings"
)

// Get returns value stored at a given path in a given decoded JSON object.
// For ex.
//   jsonStr := `{
//      "servers":[
//        {
//          "name": "server1",
//          "ip": "10.0.1.1"
//        },
//       {
//          "name": "server2"
//          "ip": "10.0.1.2"
//       }
//      ]
//   }`
//   var jsn interface{}
//   json.Unmarshal(jsonStr, &jsn)
//   firstServerName, err := jsonpath.Get(jsn, "servers[0].name")
//  Now, typecasting firstServerName.(string) will give us "server1".
//  Look at jsonpath_test.go for more examples.
func Get(obj interface{}, path string) (interface{}, error) {
  if path == "" {
    return nil, fmt.Errorf("path empty")
  }
  jsn, ok := obj.(map[string]interface{})
  if !ok {
    return nil, fmt.Errorf("error finding path %s", path)
  }
  i := strings.IndexAny(path, "[.")
  if i < 0 {
    // at leaf node in the path, extract the val and return
    val, ok := jsn[path]
    if ok {
      return val, nil
    }
    return nil, fmt.Errorf("error finding path %s", path)
  }
  // we have encountered either a . or an opening bracket [ at index i
  if path[i] == '.' {
    // TODO (sunil): Add support for escaping . character in path
    first, rest := path[:i], path[i+1:]
    return Get(jsn[first], rest)
  } else {
    // encountered an array type object, extract slice index first
    j := strings.Index(path[i+1:], "]")
    if j < 0 {
      return nil, fmt.Errorf("error closing bracket missing")
    }
    // so slice index is (i+1, i + 1 + j]
    n, err := strconv.ParseInt(path[i+1:i+1+j], 10, 64)
    if err != nil {
      return nil, err
    }
    first := path[:i]
    items, ok := jsn[first].([]interface{})
    if !ok {
      return nil, fmt.Errorf("error did not find expected slice")
    }
    if int(n) > len(items)-1 {
      return nil, fmt.Errorf("error index out of bound")
    }
    if len(path) == i+j+2 {
      return items[n], nil
    }
    rest := path[i+j+3:]
    return Get(items[n], rest)
  }
}
