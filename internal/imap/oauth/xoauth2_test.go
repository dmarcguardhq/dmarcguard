package oauth

import "testing"

func TestXOAUTH2Client_StartFormat(t *testing.T) {
	mech, ir, err := NewXOAUTH2Client("user@example.com", "tok-abc").Start()
	if err != nil {
		t.Fatal(err)
	}
	if mech != XOAUTH2 {
		t.Fatalf("got mech %q want XOAUTH2", mech)
	}
	want := "user=user@example.com\x01auth=Bearer tok-abc\x01\x01"
	if string(ir) != want {
		t.Fatalf("ir mismatch\n got: %q\nwant: %q", string(ir), want)
	}
}

func TestXOAUTH2Client_NextSurfacesChallenge(t *testing.T) {
	c := NewXOAUTH2Client("u", "t")
	resp, err := c.Next([]byte(`{"status":"401"}`))
	if err == nil {
		t.Fatal("expected error on server challenge")
	}
	if string(resp) != "" {
		t.Fatalf("expected empty ack, got %q", resp)
	}
}
