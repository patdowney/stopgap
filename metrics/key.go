package metrics

import (
	"fmt"
	"strings"
)

type Key interface {
	Add(interface{}) Key
	String() string
}

type GraphiteKey struct {
	Key string
}

func cleanse(key string) string {
	noDots := strings.Replace(key, ".", "_", -1)
	noSpaces := strings.Replace(noDots, " ", "_", -1)
	noUpper := strings.ToLower(noSpaces)
	return noUpper
}

//func (k Key) Add(newKeyPart string) Key {
func (k GraphiteKey) Add(newKeyPart interface{}) Key {
	partAsString := fmt.Sprintf("%v", newKeyPart)
	newKey := cleanse(partAsString)
	if k.Key != "" {
		newKey = fmt.Sprintf("%v.%v", k.Key, newKey)
	}

	return GraphiteKey{Key: newKey}
}

func (k GraphiteKey) String() string {
	return k.Key
}
