package parser

import (
	"bytes"
	"compress/gzip"
	"runtime"
	"testing"
)

// A DMARC report arrives from an untrusted mail server. A tiny gzip member that
// decompresses to gigabytes must be rejected, not turned into an unbounded
// allocation.
func TestGzipBombIsRejected(t *testing.T) {
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	// 512 MB of zeros -> a few hundred KB compressed, over the 50 MB cap.
	chunk := make([]byte, 1<<20)
	for i := 0; i < 512; i++ {
		_, _ = zw.Write(chunk)
	}
	_ = zw.Close()
	bomb := buf.Bytes()
	t.Logf("compressed bomb: %d bytes -> 512 MB decompressed", len(bomb))

	var m0 runtime.MemStats
	runtime.ReadMemStats(&m0)

	_, err := ParseReport(bomb)

	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)
	grewMB := (int64(m1.TotalAlloc) - int64(m0.TotalAlloc)) / (1 << 20)
	t.Logf("ParseReport err=%v, TotalAlloc grew ~%d MB", err, grewMB)

	if err == nil {
		t.Fatal("expected an error for an oversized (bomb) report, got nil")
	}
	if grewMB > 128 {
		t.Fatalf("ParseReport allocated ~%d MB despite the cap (bomb not bounded)", grewMB)
	}
}

// A normal gzipped report must still parse.
func TestValidGzippedReportStillParses(t *testing.T) {
	const xmlReport = `<?xml version="1.0"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>example.com</org_name>
    <report_id>abc123</report_id>
  </report_metadata>
  <record></record>
</feedback>`
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write([]byte(xmlReport))
	_ = zw.Close()

	fb, err := ParseReport(buf.Bytes())
	if err != nil {
		t.Fatalf("valid gzipped report should parse, got %v", err)
	}
	if fb.ReportMetadata.OrgName != "example.com" {
		t.Fatalf("unexpected org name: %q", fb.ReportMetadata.OrgName)
	}
}
