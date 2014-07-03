package metrics

import (
	"encoding/json"
	"fmt"
)

func FlattenList(key Key, value []map[string]interface{}, keyID string) map[string]string {
	pair := make(map[string]string)
	for i, v := range value {
		newKey := key.Add(fmt.Sprintf("%v", i))
		if keyID != "" {
			newKey = key.Add(fmt.Sprintf("%v", v[keyID]))
		}
		p := FlattenMap(newKey, v)
		for k, nv := range p {
			pair[k] = nv
		}
	}
	return pair
}

func FlattenPair(key Key, value interface{}) map[string]string {
	var pair map[string]string
	switch value.(type) {
	case []map[string]interface{}:
		pair = FlattenList(key, value.([]map[string]interface{}), "")
	case map[string]interface{}:
		pair = FlattenMap(key, value.(map[string]interface{}))
	case json.Number:
		pair = make(map[string]string)
		pair[key.String()] = value.(json.Number).String()
	case float64:
		pair = make(map[string]string)
		pair[key.String()] = fmt.Sprintf("%f", value)
	case bool, string, nil, []interface{}:
		break
	}
	return pair
}

func FlattenMap(prefixKey Key, m map[string]interface{}) map[string]string {
	aggregate := make(map[string]string)
	for key, value := range m {
		dotKey := prefixKey.Add(key)
		for k, v := range FlattenPair(dotKey, value) {
			aggregate[k] = v
		}
	}
	return aggregate
}
