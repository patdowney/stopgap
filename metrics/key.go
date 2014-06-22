package metrics

import (
	"fmt"
	"strings"
)

type Key struct {
	Key string
}

func cleanse(key string) string {
	return strings.ToLower(strings.Replace(key, ".", "_", -1))
}

func (k Key) Add(newKeyPart string) Key {
	newKey := cleanse(newKeyPart)
	if k.Key != "" {
		newKey = fmt.Sprintf("%v.%v", k.Key, newKey)
	}

	return Key{Key: newKey}
}

func (k Key) String() string {
	return k.Key
}
