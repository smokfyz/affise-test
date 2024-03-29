package id

import (
	"testing"
)

func TestNew(t *testing.T) {
	id := New()
	if len(id) != idSize*2 {
		t.Errorf("Expected ID length of %d, but got %d", idSize*2, len(id))
	}
}

func TestNewUnique(t *testing.T) {
	id1 := New()
	id2 := New()
	if id1 == id2 {
		t.Errorf("Expected unique IDs, but got %s and %s", id1, id2)
	}
}
