package validate

import (
	"fmt"
	"reflect"
	"testing"
)

type User struct {
	FirstName string
	Email     string
}

func TestMin(t *testing.T) {
	data := struct {
		Name string
	}{Name: "123"}
	errors, ok := Validate(data, Fields{
		"Name": Rules(Min(5)),
	})
	fmt.Println(errors)
	assertTrue(t, ok)
	asserteq(t, 0, len(errors))
}

func TestMax(t *testing.T) {
	data := struct {
		Name string
	}{Name: "12345000"}
	errors, ok := Validate(data, Fields{
		"Name": Rules(Max(4)),
	})
	assertTrue(t, ok)
	asserteq(t, 0, len(errors))
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
