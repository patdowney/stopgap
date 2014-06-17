package metrics

import "fmt"

type Key struct {
	Key string
}

func (k Key) Add(newKeyPart string) Key {
	newKey := newKeyPart
	if k.Key != "" {
		newKey = fmt.Sprintf("%v.%v", k.Key, newKeyPart)
	}

	return Key{Key: newKey}
}

func (k Key) String() string {
	return k.Key
}
