package metrics

import (
	"reflect"
	"testing"
)

func TestSimpleKeyTransform(t *testing.T) {

	tr := NewSimpleTransformer("test", "success")
	expected := map[string]interface{}{"success": nil, "something.success": nil, "something.success.somethingelse": nil}
	o := map[string]interface{}{"test": nil, "something.test": nil, "something.test.somethingelse": nil}

	result := tr.Transform(o)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("unexpected transform result: (actual:%v) != (expected:%v)", result, expected)
		t.Fail()
	}
}

func TransformRegexTransform(t *testing.T) {

	// want to transform
	// asdf.0.name = cat
	// asdf.0.stats.sometstat = 123

	// into
	// asdf.cat.name
	// asdf.cat.stats.somestat = 123

	//"(\d+)", "$1.name":
}
