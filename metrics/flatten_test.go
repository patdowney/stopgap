package metrics

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestFlattenNotImplementedPairs(t *testing.T) {
	k := Key{"key"}

	// flatten bool
	b := FlattenPair(k, false)
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.FailNow()
	}

	// flatten string,
	b = FlattenPair(k, "a string")
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.FailNow()
	}

	// flatten slice
	f := []string{"one", "two", "three"}
	b = FlattenPair(k, f)
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.Fail()
	}

	// flatten nil
	b = FlattenPair(k, nil)
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.Fail()
	}
}

func TestFlattenPairs(t *testing.T) {
	k := Key{"key"}

	// flatten json.Number - float
	b := FlattenPair(k, json.Number("1.23"))
	if b == nil {
		t.Errorf("failed to flatten a json.Number: %v", b)
		t.FailNow()
	}
	if b[k.String()] != "1.23" {
		t.Errorf("failed to flatten a json.Number to float correctly: %v != %v", b, "1.23")
		t.FailNow()
	}

	// flatten json.Number - int
	b = FlattenPair(k, json.Number("1"))
	if b == nil {
		t.Errorf("failed to flatten a json.Number: %v", b)
		t.FailNow()
	}
	if b[k.String()] != "1" {
		t.Errorf("failed to flatten a json.Number to int correctly: %v != %v", b[k.String()], "1")
		t.FailNow()
	}

	// flatten float64
	b = FlattenPair(k, float64(1.234))
	if b == nil {
		t.Errorf("failed to flatten a float64: %v", b)
		t.Fail()
	}
	if b[k.String()] != "1.234000" {
		t.Errorf("failed to flatten a json.Number to int correctly: %v != %v", b[k.String()], 1.234000)
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

func TestFlattenListKey(t *testing.T) {
	k := Key{}
	expectedKey := "0.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"number": 2.0}}

	b := FlattenList(k, input, "")
	if b == nil {
		t.Errorf("failed to flatten array: %v", b)
		t.Fail()
	}

	keys := mapKeys(b)
	if keys[0] != expectedKey {
		t.Errorf("failed to flatten array key properly: %v != %v", keys[0], expectedKey)
		t.FailNow()
	}
}

func TestFlattenListCustomKey(t *testing.T) {
	k := Key{}
	expectedKey := "test.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"name":   "test",
			"number": 2.0}}

	b := FlattenList(k, input, "name")
	if b == nil {
		t.Errorf("failed to flatten array: %v", b)
		t.Fail()
	}

	keys := mapKeys(b)
	if keys[0] != expectedKey {
		t.Errorf("failed to flatten array key properly: %v != %v", keys[0], expectedKey)
		t.FailNow()
	}
}
