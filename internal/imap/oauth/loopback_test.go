package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

func TestGeneratePKCE_VerifierAndChallengeMatch(t *testing.T) {
	verifier, challenge, err := generatePKCE()
	if err != nil {
		t.Fatal(err)
	}
	if len(verifier) < 43 || len(verifier) > 128 {
		t.Fatalf("verifier length %d outside RFC 7636 bounds (43-128)", len(verifier))
	}
	sum := sha256.Sum256([]byte(verifier))
	want := base64.RawURLEncoding.EncodeToString(sum[:])
	if challenge != want {
		t.Fatal("challenge does not match S256(verifier)")
	}
}

func TestRandomToken_Distinct(t *testing.T) {
	a, err := randomToken(32)
	if err != nil {
		t.Fatal(err)
	}
	b, err := randomToken(32)
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("two random tokens collided — RNG broken")
	}
}
