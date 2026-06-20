package auth

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func newSigner(t *testing.T) *SessionSigner {
	t.Helper()
	// 64-char hex string — not valid base64, falls through to raw 64-byte key.
	s, err := NewSessionSigner("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef!", 0)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestSessionSigner_RoundTrip(t *testing.T) {
	s := newSigner(t)
	now := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	cookie, err := s.Sign("alice", "alice@example.com", now)
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Verify(cookie, now.Add(time.Hour))
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if got.Subject != "alice" || got.Email != "alice@example.com" {
		t.Fatalf("unexpected session: %+v", got)
	}
}

func TestSessionSigner_TamperDetected(t *testing.T) {
	s := newSigner(t)
	cookie, _ := s.Sign("alice", "a@example.com", time.Now())
	// Flip a byte in the payload portion.
	tampered := "X" + cookie[1:]
	if _, err := s.Verify(tampered, time.Now()); !errors.Is(err, ErrSessionBadSig) && !errors.Is(err, ErrSessionMalformed) {
		t.Fatalf("expected sig/malformed error, got %v", err)
	}
}

func TestSessionSigner_ExpirationEnforced(t *testing.T) {
	s, err := NewSessionSigner("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef!", 1)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	cookie, _ := s.Sign("a", "a@x.com", now)
	if _, err := s.Verify(cookie, now.Add(48*time.Hour)); !errors.Is(err, ErrSessionExpired) {
		t.Fatalf("expected ErrSessionExpired, got %v", err)
	}
}

func TestSessionSigner_RejectsShortSecret(t *testing.T) {
	if _, err := NewSessionSigner("short", 0); err == nil {
		t.Fatal("expected error on short secret")
	}
}

func TestSessionSigner_StateRoundTrip(t *testing.T) {
	s := newSigner(t)
	now := time.Now()
	state, err := s.SignState(now)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.VerifyState(state, now.Add(time.Minute)); err != nil {
		t.Fatalf("VerifyState: %v", err)
	}
}

func TestSessionSigner_StateExpires(t *testing.T) {
	s := newSigner(t)
	now := time.Now()
	state, _ := s.SignState(now)
	if err := s.VerifyState(state, now.Add(time.Hour)); !errors.Is(err, ErrSessionExpired) {
		t.Fatalf("expected ErrSessionExpired, got %v", err)
	}
}

func TestGenerateSessionSecret_DistinctAndLongEnough(t *testing.T) {
	a, _ := GenerateSessionSecret()
	b, _ := GenerateSessionSecret()
	if a == b {
		t.Fatal("two generated secrets collided")
	}
	// base64 of 32 bytes is 44 chars including padding.
	if len(a) < 40 || strings.ContainsAny(a, " \t\n") {
		t.Fatalf("unexpected secret format: %q", a)
	}
}
