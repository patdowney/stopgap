package metrics

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestDecodeNotImplementedPairs(t *testing.T) {
	k := Key{"key"}

	// decode bool
	b := decodePair(k, false)
	if b != nil {
		t.Errorf("decoded unexpected value: %v", b)
		t.FailNow()
	}

	// decode string,
	b = decodePair(k, "a string")
	if b != nil {
		t.Errorf("decoded unexpected value: %v", b)
		t.FailNow()
	}

	// decode slice
	f := []string{"one", "two", "three"}
	b = decodePair(k, f)
	if b != nil {
		t.Errorf("decoded unexpected value: %v", b)
		t.Fail()
	}

	// decode nil
	b = decodePair(k, nil)
	if b != nil {
		t.Errorf("decoded unexpected value: %v", b)
		t.Fail()
	}
}

func TestDecodePairs(t *testing.T) {
	k := Key{"key"}

	// decode json.Number - float
	b := decodePair(k, json.Number("1.23"))
	if b == nil {
		t.Errorf("failed to decode a json.Number: %v", b)
		t.FailNow()
	}
	if b[k.String()] != "1.23" {
		t.Errorf("failed to decode a json.Number to float correctly: %v != %v", b, "1.23")
		t.FailNow()
	}

	// decode json.Number - int
	b = decodePair(k, json.Number("1"))
	if b == nil {
		t.Errorf("failed to decode a json.Number: %v", b)
		t.FailNow()
	}
	if b[k.String()] != "1" {
		t.Errorf("failed to decode a json.Number to int correctly: %v != %v", b[k.String()], "1")
		t.FailNow()
	}

	// decode float64
	b = decodePair(k, float64(1.234))
	if b == nil {
		t.Errorf("failed to decode a float64: %v", b)
		t.Fail()
	}
	if b[k.String()] != "1.234000" {
		t.Errorf("failed to decode a json.Number to int correctly: %v != %v", b[k.String()], 1.234000)
		t.FailNow()
	}
}

func mapKeys(m map[string]string) []string {
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestMapKeys(t *testing.T) {
	expectedKeys := []string{"key1", "key2", "key3"}
	m := map[string]string{
		"key1": "", "key2": "", "key3": ""}

	keys := mapKeys(m)
	if keys[0] != expectedKeys[0] ||
		keys[1] != expectedKeys[1] ||
		keys[2] != expectedKeys[2] {
		t.Errorf("failed to list correct keys: %v != %v", keys, expectedKeys)
		t.FailNow()
	}
}

func TestDecodeListKey(t *testing.T) {
	k := Key{}
	expectedKey := "0.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"number": 2.0}}

	b := decodeList(k, input, "")
	if b == nil {
		t.Errorf("failed to decode array: %v", b)
		t.Fail()
	}

	keys := mapKeys(b)
	if keys[0] != expectedKey {
		t.Errorf("failed to decode array key properly: %v != %v", keys[0], expectedKey)
		t.FailNow()
	}
}

func TestDecodeListCustomKey(t *testing.T) {
	k := Key{}
	expectedKey := "test.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"name":   "test",
			"number": 2.0}}

	b := decodeList(k, input, "name")
	if b == nil {
		t.Errorf("failed to decode array: %v", b)
		t.Fail()
	}

	keys := mapKeys(b)
	if keys[0] != expectedKey {
		t.Errorf("failed to decode array key properly: %v != %v", keys[0], expectedKey)
		t.FailNow()
	}
}

func TestDecodeDictKeys(t *testing.T) {
	k := Key{}
	expectedKey := "number"

	input := map[string]interface{}{
		"number": 2.0}

	b := decodeDict(k, input)
	if b == nil {
		t.Errorf("failed to decode dict: %v", b)
		t.Fail()
	}

	keys := mapKeys(b)
	if keys[0] != expectedKey {
		t.Errorf("failed to decode array key properly: %v != %v", keys[0], expectedKey)
		t.FailNow()
	}
}
