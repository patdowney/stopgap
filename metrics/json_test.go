package metrics

import (
	"encoding/json"
	"testing"
)

func TestDecodePair(t *testing.T) {
	k := Key{"key"}
	// decodePair(key Key, value interface{}) map[string]string

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

	// decode slice - blows up?
	//f := []string{"one", "two", "three"}
	//b = decodePair(k, f) //[]interface{"one", "two", "three"})
	//if b != nil {
	//		t.Fail()
	//	}

	// decode json.Number - float
	b = decodePair(k, json.Number("1.23"))
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

	// decode int
	//b = decodePair(k, 1)
	// decode map
}

func TestDecodeDict(t *testing.T) {
	// decodeDict(prefixKey Key, dict map[string]interface{}) map[string]string {
}

func TestDecode(t *testing.T) {
	// (d *JSONMetricDecoder) Decode(pairs *[]Metric) error {
}

func TestDecoder(t *testing.T) {
	// NewDecoder(reader io.Reader) *JSONMetricDecoder {
}
