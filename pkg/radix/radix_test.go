package radix

import (
	"testing"
)

func TestSet(t *testing.T) {

	tree := NewTree()

	tree.Set("abc", "1")
	tree.Set("def", "2")
	tree.Set("abd", "3")
}

func TestLen(t *testing.T) {

	tree := NewTree()

	tree.Set("abc", "1")
	tree.Set("def", "2")
	tree.Set("abd", "3")

	s := tree.Size()
	if s != 3 {
		t.Fatalf("tree.Size() failed, expected %d, have %d", 3, s)
	}
}

func TestGet(t *testing.T) {

	tree := NewTree()

	tree.Set("abc", "1")
	tree.Set("def", "2")
	tree.Set("abd", "3")

	if v, found := tree.Get("def"); !found {
		t.Fatalf("tree.Get() failed, no results returned")
	} else {
		if v != "2" {
			t.Fatalf("tree.Get() failed, expected %v, have %v", "2", v)
		}
	}

	if v, found := tree.Get("abd"); !found {
		t.Fatalf("tree.Get() failed, no results returned")
	} else {
		if v != "3" {
			t.Fatalf("tree.Get() failed, expected %v, have %v", "3", v)
		}
	}

	if v, found := tree.Get("ab"); found {
		t.Fatalf("tree.Get() failed, no results expected, have %v", v)
	}

}

func TestDelete(t *testing.T) {

	tree := NewTree()

	tree.Set("a", "1")
	tree.Set("ab", "2")
	tree.Set("acd", "3")
	tree.Set("acdx", "4")

	err := tree.Delete("acd")
	if err != nil {
		t.Fatalf("tree.Delete() failed: %v", err)
	}
	if s := tree.Size(); s != 3 {
		t.Fatalf("tree.Delete() failed, expected %d, have %d", 3, s)
	}
}

func TestGetClosest(t *testing.T) {

	tree := NewTree()

	tree.Set("a", "1")
	tree.Set("ab", "2")
	tree.Set("acd", "3")
	tree.Set("acdx", "4")

	longest, value, found := tree.GetClosest("acde")
	if !found {
		t.Fatalf("tree.GetClosest() failed: no results found")
	}
	if longest != "acd" {
		t.Fatalf("tree.GetClosest() failed, expected longest prefix to be %v, have %v", "acd", longest)
	}
	if value != "3" {
		t.Fatalf("tree.GetClosest() failed, expected closest value to be %v, have %v", "3", value)
	}
}

func TestString(t *testing.T) {
	tree := NewTree()

	tree.Set("a", "1")
	tree.Set("ab", "2")
	tree.Set("acd", "3")
	tree.Set("acdx", "4")

	s := tree.String()

	t.Log(s)
}
