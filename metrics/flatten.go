package metrics

import "encoding/json"

func Flatten(key Key, value interface{}, keyID string) map[string]interface{} {
	return FlattenPair(key, value, keyID)
}

func FlattenMapList(key Key, mapList []map[string]interface{}, keyID string) map[string]interface{} {
	pair := make(map[string]interface{})
	for i, v := range mapList {
		newKey := key.Add(i)
		if keyVal, ok := v[keyID]; ok {
			newKey = key.Add(keyVal)
		}
		//p := FlattenMap(newKey, v, keyID)
		p := Flatten(newKey, v, keyID)
		for k, nv := range p {
			pair[k] = nv
		}
	}
	return pair
}

func flattenJSONNumber(key Key, value json.Number) map[string]interface{} {
	pair := make(map[string]interface{})

	intValue, err := value.Int64()
	if err != nil {
		floatValue, _ := value.Float64()
		pair[key.String()] = floatValue
	} else {
		pair[key.String()] = intValue
	}
	return pair
}

func FlattenPair(key Key, value interface{}, keyID string) map[string]interface{} {
	var pair map[string]interface{}
	switch value.(type) {
	case []map[string]interface{}:
		return FlattenMapList(
			key,
			value.([]map[string]interface{}),
			keyID)
	case map[string]interface{}:
		return FlattenMap(
			key,
			value.(map[string]interface{}),
			keyID)
	case json.Number:
		return flattenJSONNumber(key, value.(json.Number))
	case float64, bool, string:
		return map[string]interface{}{key.String(): value}
	case nil, []interface{}:
		break
	}
	return pair
}

func FlattenMap(prefixKey Key, m map[string]interface{}, keyID string) map[string]interface{} {
	aggregate := make(map[string]interface{})
	for key, value := range m {
		dotKey := prefixKey.Add(key)
		for k, v := range FlattenPair(dotKey, value, keyID) {
			aggregate[k] = v
		}
	}
	return aggregate
}
