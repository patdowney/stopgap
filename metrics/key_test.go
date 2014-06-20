package metrics

import "testing"

func TestAddFirstKey(t *testing.T) {
	empty := Key{}
	expected := "first"

	k := empty.Add("first")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}
}

func TestAddMultipleKeys(t *testing.T) {
	empty := Key{}
	expected := "first.second.third"

	k := empty.Add("first").Add("second").Add("third")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}
}
