package metrics

import "testing"

func TestDecodePair(t *testing.T) {
	// decodePair(key Key, value interface{}) map[string]string
	// decode bool
	// decode string,
	// decode slice
	//decode json.Number
	// decode flaot64
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
