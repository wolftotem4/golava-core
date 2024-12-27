package auth

import "testing"

func TestMatchID(t *testing.T) {
	recaller := NewRecaller(1, "token", "hash")
	if !recaller.MatchID(1) {
		t.Error("MatchID failed")
	}

	if recaller.MatchID(2) {
		t.Error("MatchID failed")
	}

	if !recaller.MatchID("1") {
		t.Error("MatchID failed")
	}

	if recaller.MatchID("2") {
		t.Error("MatchID failed")
	}

	recaller = NewRecaller("1", "token", "hash")
	if !recaller.MatchID(1) {
		t.Error("MatchID failed")
	}

	if recaller.MatchID(2) {
		t.Error("MatchID failed")
	}

	if !recaller.MatchID("1") {
		t.Error("MatchID failed")
	}

	if recaller.MatchID("2") {
		t.Error("MatchID failed")
	}
}
