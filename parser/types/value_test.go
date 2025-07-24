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
	if !tuple.CanContain(new(List)) {
		t.Error("Tuple must be able to contain anything")
	}
}

func TestList(t *testing.T) {
	list := new(List)
	list.AddValues(Int(1), Int(2), Int(3))
	if !list.Type().Is(NewListType(IntType)) {
		t.Error("Expected int, got", list.Type())
	}
	val := list.Value().([]Value)
	for i, value := range list.GetValues() {
		if val[i] != value {
			t.Error("Expected", value, "got", val[i])
		}
	}
	if len(val) != 3 {
		t.Error("Expected 3, got", len(val))
	}
	if !val[0].Type().Is(IntType) {
		t.Error("Expected string, got", val[0].Type())
	}
	v, ok := val[0].Value().(int)
	if !ok {
		t.Error("Cannot convert value to string")
	} else {
		if v != 1 {
			t.Error("Expected 1, got", v)
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
	if !val[2].Type().Is(IntType) {
		t.Error("Expected int, got", val[2].Type())
	}
	v, ok = val[2].Value().(int)
	if !ok {
		t.Error("Cannot convert value to string")
	} else {
		if v != 3 {
			t.Error("Expected 3, got", v)
		}
	}
	if !list.CanContain(Int(0)) {
		t.Error("List must be able to contain the same type")
	}
	if list.CanContain(String("A")) {
		t.Error("List must contain only values with the same type")
	}
}
