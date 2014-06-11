package lbenv

import (
	"reflect"
	"testing"
)

func TestValues(t *testing.T) {
	env := New()
	if env.Has("MY_PATH") {
		t.Fatalf("expected no 'MY_PATH' env.var")
	}
	err := env.Append("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("problem appending: %v", err)
	}

	if !env.Has("MY_PATH") {
		t.Fatalf("expected a 'MY_PATH' env.var")
	}

	err = env.Append("MY_PATH", "newValue:secondVal:valval")
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	v := splitpath(env.Get("MY_PATH").Value)

	if v[len(v)-1] != "valval" {
		t.Fatalf("expected %q. got=%q", "valval", v[len(v)-1])
	}

	if !in_str_slice("newValue", v) {
		t.Fatalf("expected %q to be in slice %v", "newValue", v)
	}

	err = env.Remove("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if in_str_slice("newValue", v) {
		t.Fatalf("expected %q NOT to be in slice %v", "newValue", v)
	}

	err = env.Prepend("MY_PATH", "newValue")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if !in_str_slice("newValue", v) {
		t.Fatalf("expected %q to be in slice %v", "newValue", v)
	}

	err = env.Set("MY_PATH", "hi:hello")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	v = splitpath(env.Get("MY_PATH").Value)
	if !reflect.DeepEqual(v, []string{"hi", "hello"}) {
		t.Fatalf("expected env=%v. got=%v", "newValue", []string{"hi", "hello"}, v)
	}

	err = env.Unset("MY_PATH")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if env.Has("MY_PATH") {
		t.Fatalf("expected no 'MY_PATH' env.var")
	}
}
