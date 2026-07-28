package ui

import "testing"

func TestDetectColorDisabled(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	if detectColor() {
		t.Fatal("expected color disabled when NO_COLOR set")
	}
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "dumb")
	if detectColor() {
		t.Fatal("expected color disabled when TERM=dumb")
	}
}
