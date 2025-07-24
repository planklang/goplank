package types

import (
	"testing"
)

func TestTuple(t *testing.T) {
	tuple := new(Tuple)
	tuple.AddValues(String("A"), Int(2), String("C"))
	if !tuple.Type().Is(NewTupleType(StringType, IntType, StringType)) {
		t.Error("Expected (string int string), got", tuple.Type())
	}
	val := tuple.Value().([]Value)
	for i, value := range tuple.GetValues() {
		if val[i] != value {
			t.Error("Expected", value, "got", val[i])
		}
	}
	if len(val) != 3 {
		t.Error("Expected 3, got", len(val))
	}
	if !val[0].Type().Is(StringType) {
		t.Error("Expected string, got", val[0].Type())
	}
	v, ok := val[0].Value().(string)
	if !ok {
		t.Error("Cannot convert value to string")
	} else {
		if v != "A" {
			t.Error("Expected A, got", v)
		}
	}
	if !val[1].Type().Is(IntType) {
		t.Error("Expected int, got", val[1].Type())
	}
	vi, ok := val[1].Value().(int)
	if !ok {
		t.Error("Cannot convert value to int")
	} else {
		if vi != 2 {
			t.Error("Expected 2, got", vi)
		}
	}
	if !val[2].Type().Is(StringType) {
		t.Error("Expected string, got", val[2].Type())
	}
	v, ok = val[2].Value().(string)
	if !ok {
		t.Error("Cannot convert value to string")
	} else {
		if v != "C" {
			t.Error("Expected C, got", v)
		}
	}
}
