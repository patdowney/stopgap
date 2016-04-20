package metrics

import (
	"regexp"
	"strings"
)

// Transformer ...
type Transformer interface {
	Transform(map[string]interface{}) map[string]interface{}
}

// SimpleKeyTransformer ...
type SimpleKeyTransformer struct {
	Match   string
	Replace string
}

// RegexpKeyTransformer ...
type RegexpKeyTransformer struct {
	Match   regexp.Regexp
	Replace regexp.Regexp
}

func (t *SimpleKeyTransformer) transformKey(key string) string {
	keyParts := strings.Split(key, ".")
	newParts := make([]string, len(keyParts))
	for i, k := range keyParts {
		if k == t.Match {
			newParts[i] = t.Replace
		} else {
			newParts[i] = k
		}
	}

	return strings.Join(newParts, ".")
}

// Transform ...
func (t *SimpleKeyTransformer) Transform(original map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})
	for k, v := range original {
		y := t.transformKey(k)
		result[y] = v
	}
	return result
}

// NewSimpleTransformer ...
func NewSimpleTransformer(match string, replace string) Transformer {
	return &SimpleKeyTransformer{
		Match:   match,
		Replace: replace}
}
