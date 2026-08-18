package imap

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/rs/zerolog"
)

const feedbackXML = `<?xml version="1.0"?><feedback><report_metadata><report_id>1</report_id></report_metadata></feedback>`

func zipped(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("google.com!example.com!1!2.xml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write([]byte(feedbackXML)); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// A Google report relayed as a message/rfc822 attachment (the ".msg" case).
func TestCollectAttachmentsNestedMessage(t *testing.T) {
	inner := "From: noreply-dmarc-support@google.com\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=in\r\n\r\n" +
		"--in\r\nContent-Type: application/zip\r\n" +
		"Content-Disposition: attachment; filename=\"report.zip\"\r\n" +
		"Content-Transfer-Encoding: base64\r\n\r\n" +
		encodeB64(zipped(t)) + "\r\n--in--\r\n"

	outer := "From: relay@example.com\r\nMIME-Version: 1.0\r\n" +
		"Content-Type: multipart/mixed; boundary=out\r\n\r\n" +
		"--out\r\nContent-Type: message/rfc822\r\n" +
		"Content-Disposition: attachment; filename=\"Report domain.msg\"\r\n\r\n" +
		inner + "\r\n--out--\r\n"

	log := zerolog.Nop()
	c := &Client{log: &log}
	atts := c.collectAttachments(strings.NewReader(outer), 0)
	if len(atts) != 1 {
		t.Fatalf("want 1 attachment, got %d", len(atts))
	}
	if !bytes.HasPrefix(atts[0].Data, []byte("PK\x03\x04")) {
		t.Fatalf("want zip payload, got %q", atts[0].Data[:4])
	}
}

func TestIsDMARCAttachmentSniff(t *testing.T) {
	if !isDMARCAttachment("opaque.bin", []byte(feedbackXML)) {
		t.Error("content sniff failed for raw feedback XML")
	}
	if isDMARCAttachment("logo.png", []byte("\x89PNG")) {
		t.Error("PNG should not be treated as a DMARC report")
	}
}

func encodeB64(b []byte) string {
	e := base64.StdEncoding.EncodeToString(b)
	var out []string
	for len(e) > 76 {
		out, e = append(out, e[:76]), e[76:]
	}
	return strings.Join(append(out, e), "\r\n")
}
