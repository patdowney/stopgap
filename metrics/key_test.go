package metrics

import "testing"

func TestGraphiteKeyAddFirst(t *testing.T) {
	empty := GraphiteKey{}
	expected := "first"

	k := empty.Add("first")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}
}

func TestGraphiteKeyAddMultiple(t *testing.T) {
	empty := GraphiteKey{}
	expected := "first.second.third"

	k := empty.Add("first").Add("second").Add("third")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}
}

func TestGraphiteKeyAddWithDots(t *testing.T) {
	empty := GraphiteKey{}
	expected := "first_nodots"

	k := empty.Add("first.nodots")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}

func TestGraphiteKeyAddWithUpperCase(t *testing.T) {
	empty := GraphiteKey{}
	expected := "shouldlowercase"

	k := empty.Add("SHouLdLowerCasE")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}

func TestGraphiteKeyAddWithSpace(t *testing.T) {
	empty := GraphiteKey{}
	expected := "should_not_have_any_spaces"

	k := empty.Add("should not have any spaces")

	if k.String() != expected {
		t.Errorf("'%v' != '%v'", expected, k)
		t.FailNow()
	}

}
