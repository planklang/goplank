package types

import "testing"

func TestLiteralType_Is(t *testing.T) {
	generalTest(t, []Type{StringType, IntType, FloatType, DefaultLiteralType})
}

func TestTupleType_Is(t *testing.T) {
	generalTest(t, []Type{NewTupleType(StringType), NewTupleType(IntType), NewTupleType(StringType, IntType)})
}

func TestListType_Is(t *testing.T) {
	generalTest(t, []Type{NewListType(StringType), NewListType(IntType), NewListType(FloatType)})
}

func generalTest(t *testing.T, typs []Type) {
	for i, typ := range typs {
		for j, t2 := range typs {
			if i == j {
				if !typ.Is(t2) {
					t.Errorf("%s is not a types %s", typ, t2)
				}
			} else {
				if typ.Is(t2) {
					t.Errorf("%s is not a types %s", typ, t2)
				}
			}
		}
	}
}
