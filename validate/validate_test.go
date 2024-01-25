package validate

import (
	"fmt"
	"reflect"
	"testing"
)

type CustomError struct {
	FirstName string
	Name      string
}

func TestBarBar(t *testing.T) {
	data := struct {
		FirstName string
		LastName  string
		Email     string
	}{
		FirstName: "",
		LastName:  "GG",
		Email:     "cryptoanthdm@gmail.com",
	}

	var errs CustomError
	ok := New(data, Fields{
		"FirstName": Rules(
			Required,
			Min(10),
			Max(100),
			Message("The name needs to present")),
	}).Validate(errs)
	if !ok {
		fmt.Println(errs)
	}
}

func TestCustomMessage(t *testing.T) {
	data := struct {
		Name string
	}{Name: ""}
	errs := map[string]string{}
	ok := New(data, Fields{
		"Name": Rules(Required, Message("name not good")),
	}).Validate(errs)
	assertFalse(t, ok)
	asserteq(t, 1, len(errs))
	asserteq(t, "name not good", errs["Name"])
}

func TestRequired(t *testing.T) {
	data := struct {
		Name string
	}{Name: ""}
	t.Run("invalid", func(t *testing.T) {
		errs := map[string]string{}
		ok := New(data, Fields{
			"Name": Rules(Required),
		}).Validate(errs)
		assertFalse(t, ok)
		asserteq(t, 1, len(errs))
	})
	t.Run("valid", func(t *testing.T) {
		errs := map[string]string{}
		data.Name = "foo"
		ok := New(data, Fields{
			"Name": Rules(Required),
		}).Validate(errs)
		assertTrue(t, ok)
		asserteq(t, 0, len(errs))
	})
}

func TestUrl(t *testing.T) {
	data := struct {
		Url string
	}{Url: "http://foocom"}
	t.Run("invalid", func(t *testing.T) {
		errs := map[string]string{}
		ok := New(data, Fields{
			"Url": Rules(Url),
		}).Validate(errs)
		assertFalse(t, ok)
		asserteq(t, 1, len(errs))
	})
	t.Run("valid", func(t *testing.T) {
		errs := map[string]string{}
		data.Url = "http://foo.com"
		ok := New(data, Fields{
			"Url": Rules(Url),
		}).Validate(errs)
		assertTrue(t, ok)
		asserteq(t, 0, len(errs))
	})
}

func TestEmail(t *testing.T) {
	data := struct {
		Email string
	}{Email: "agg.com"}
	t.Run("invalid", func(t *testing.T) {
		errs := map[string]string{}
		ok := New(data, Fields{
			"Email": Rules(Email),
		}).Validate(errs)
		assertFalse(t, ok)
		asserteq(t, 1, len(errs))
	})
	t.Run("valid", func(t *testing.T) {
		errs := map[string]string{}
		data.Email = "a@gg.com"
		ok := New(data, Fields{
			"Email": Rules(Email),
		}).Validate(errs)
		assertTrue(t, ok)
		asserteq(t, 0, len(errs))
	})
}

func TestMin(t *testing.T) {
	data := struct {
		Name string
	}{Name: "123"}
	t.Run("invalid", func(t *testing.T) {
		errs := map[string]string{}
		ok := New(data, Fields{
			"Name": Rules(Min(5)),
		}).Validate(errs)
		assertFalse(t, ok)
		asserteq(t, 1, len(errs))
	})
	t.Run("valid", func(t *testing.T) {
		errs := map[string]string{}
		data.Name = "123456"
		ok := New(data, Fields{
			"Name": Rules(Min(5)),
		}).Validate(errs)
		assertTrue(t, ok)
		asserteq(t, 0, len(errs))
	})
}

func TestMax(t *testing.T) {
	data := struct {
		Name string
	}{Name: "1234444444"}
	t.Run("invalid", func(t *testing.T) {
		errs := map[string]string{}
		ok := New(data, Fields{
			"Name": Rules(Max(5)),
		}).Validate(errs)
		assertFalse(t, ok)
		asserteq(t, 1, len(errs))
	})
	t.Run("valid", func(t *testing.T) {
		errs := map[string]string{}
		foo := struct {
			Name string
		}{Name: "123"}
		ok := New(foo, Fields{
			"Name": Rules(Max(5)),
		}).Validate(errs)
		assertTrue(t, ok)
		asserteq(t, 0, len(errs))
	})
}

func assertTrue(t *testing.T, con bool) {
	if !con {
		t.Fatalf("expected true")
	}
}

func assertFalse(t *testing.T, con bool) {
	if con {
		t.Fatalf("expected false")
	}
}

func asserteq(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("expected %v to equal %v", a, b)
	}
}
