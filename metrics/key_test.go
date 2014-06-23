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

func TestAddKeyWithDots(t *testing.T) {
	empty := Key{}
	expected := "first_nodots"

	k := empty.Add("first.nodots")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}

func TestAddKeyWithUpperCase(t *testing.T) {
	empty := Key{}
	expected := "shouldlowercase"

	k := empty.Add("SHouLdLowerCasE")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}

func TestAddKeyWithSpace(t *testing.T) {
	empty := Key{}
	expected := "should_not_have_any_spaces"

	k := empty.Add("should not have any spaces")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}
