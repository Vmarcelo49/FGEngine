package types

import (
	"testing"
)

func TestBoxTypeNames(t *testing.T) {
	if Collision != "collision" || Hit != "hit" || Hurt != "hurt" {
		t.Fatalf("wire names changed: %q %q %q", Collision, Hit, Hurt)
	}
	if len(BoxTypes) != 3 {
		t.Fatalf("BoxTypes should list exactly 3 types, got %d", len(BoxTypes))
	}
	if BoxTypes[0] != Collision || BoxTypes[1] != Hit || BoxTypes[2] != Hurt {
		t.Fatalf("BoxTypes order changed: %v", BoxTypes)
	}
	if Collision.String() != "Collision" || BoxType("nope").String() != "Unknown" {
		t.Fatal("String() mismatch")
	}
}
