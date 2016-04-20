package metrics

import (
	"fmt"
	"strings"
)

// Key ...
type Key interface {
	Add(interface{}) Key
	String() string
}

// GraphiteKey ...
type GraphiteKey struct {
	Key string
}

func cleanse(key string) string {
	noDots := strings.Replace(key, ".", "_", -1)
	noSpaces := strings.Replace(noDots, " ", "_", -1)
	noUpper := strings.ToLower(noSpaces)
	return noUpper
}

// Add ...
func (k GraphiteKey) Add(newKeyPart interface{}) Key {
	partAsString := fmt.Sprintf("%v", newKeyPart)
	newKey := cleanse(partAsString)
	if k.Key != "" {
		newKey = fmt.Sprintf("%v.%v", k.Key, newKey)
	}

	return GraphiteKey{Key: newKey}
}

// String ...
func (k GraphiteKey) String() string {
	return k.Key
}
