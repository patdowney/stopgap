package metrics

import (
	"encoding/json"
	"sort"
	"testing"
)

func TestFlattenNotImplementedPairs(t *testing.T) {
	k := GraphiteKey{"key"}

	// flatten slice
	f := []string{"one", "two", "three"}
	b := FlattenPair(k, f, "")
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.Fail()
	}

	// flatten nil
	b = FlattenPair(k, nil, "")
	if b != nil {
		t.Errorf("flatten unexpected value: %v", b)
		t.Fail()
	}
}

func testFlattenPair(t *testing.T, value interface{}, expected interface{}, label string) {
	k := GraphiteKey{"key"}
	b := FlattenPair(k, value, "")
	if b == nil {
		t.Errorf("(%v) failed to flatten pair: %v", label, b)
		t.FailNow()
	}
	if b[k.String()] != expected {
		t.Errorf("(%v) failed to flatten correctly: %v != %v", label, b[k.String()], value)
		t.FailNow()
	}
}

func TestFlattenPairs(t *testing.T) {
	testFlattenPair(t, false, false, "bool")

	testFlattenPair(t, "a string", "a string", "string")

	testFlattenPair(t, json.Number("1.23"), 1.23, "json.Number(float)")

	testFlattenPair(t, json.Number("1"), int64(1), "json.Number(int)")

	testFlattenPair(t, float64(1.234), 1.234, "float64")
}

func mapKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
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
	k := GraphiteKey{}
	expectedKey := "0.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"number": 2.0}}

	b := FilterMapNumbers(FlattenMapList(k, input, ""))
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
	k := GraphiteKey{}
	expectedKey := "test.number"

	input := []map[string]interface{}{
		map[string]interface{}{
			"name":   "test",
			"number": 2.0}}

	b := FilterMapNumbers(FlattenMapList(k, input, "name"))
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
