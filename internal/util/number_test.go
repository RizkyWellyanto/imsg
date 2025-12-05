package util

import "testing"

func TestNormalizeE164(t *testing.T) {
	got := NormalizeE164("(415) 555-1212", "US")
	want := "+14155551212"
	if got != want {
		t.Fatalf("want %s got %s", want, got)
	}
}

func TestNormalizeE164Fallback(t *testing.T) {
	raw := "not-a-number"
	if got := NormalizeE164(raw, "US"); got != raw {
		t.Fatalf("expected fallback to raw, got %s", got)
	}
}
