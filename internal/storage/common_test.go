package storage

import (
	"fmt"
	"strings"
	"testing"

	"github.com/goccy/go-json"

	"github.com/meysam81/parse-dmarc/internal/parser"
)

func TestGetStatistics_HasData(t *testing.T) {
	// Create an in-memory SQLite database for testing
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	t.Run("empty database", func(t *testing.T) {
		stats, err := storage.GetStatistics()
		if err != nil {
			t.Fatalf("Failed to get statistics: %v", err)
		}

		if stats.HasData {
			t.Errorf("Expected HasData to be false for empty database, got true")
		}

		if stats.TotalReports != 0 {
			t.Errorf("Expected TotalReports to be 0, got %d", stats.TotalReports)
		}

		if stats.TotalMessages != 0 {
			t.Errorf("Expected TotalMessages to be 0, got %d", stats.TotalMessages)
		}

		if stats.CompliantMessages != 0 {
			t.Errorf("Expected CompliantMessages to be 0, got %d", stats.CompliantMessages)
		}

		if stats.ComplianceRate != 0 {
			t.Errorf("Expected ComplianceRate to be 0, got %f", stats.ComplianceRate)
		}
	})

	t.Run("database with report", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>google.com</org_name>
    <email>noreply-dmarc-support@google.com</email>
    <report_id>12345678901234567890</report_id>
    <date_range>
      <begin>1609459200</begin>
      <end>1609545600</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>none</p>
    <sp>none</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip>192.0.2.1</source_ip>
      <count>100</count>
      <policy_evaluated>
        <disposition>none</disposition>
        <dkim>pass</dkim>
        <spf>pass</spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results>
      <spf>
        <domain>example.com</domain>
        <result>pass</result>
      </spf>
      <dkim>
        <domain>example.com</domain>
        <result>pass</result>
      </dkim>
    </auth_results>
  </record>
</feedback>`

		feedback, err := parser.ParseReport([]byte(xmlData))
		if err != nil {
			t.Fatalf("Failed to parse report: %v", err)
		}

		err = storage.SaveReport(feedback)
		if err != nil {
			t.Fatalf("Failed to save report: %v", err)
		}

		stats, err := storage.GetStatistics()
		if err != nil {
			t.Fatalf("Failed to get statistics after adding report: %v", err)
		}

		if !stats.HasData {
			t.Errorf("Expected HasData to be true after adding report, got false")
		}

		if stats.TotalReports != 1 {
			t.Errorf("Expected TotalReports to be 1, got %d", stats.TotalReports)
		}

		if stats.TotalMessages != 100 {
			t.Errorf("Expected TotalMessages to be 100, got %d", stats.TotalMessages)
		}

		if stats.CompliantMessages != 100 {
			t.Errorf("Expected CompliantMessages to be 100, got %d", stats.CompliantMessages)
		}

		if stats.ComplianceRate != 100.0 {
			t.Errorf("Expected ComplianceRate to be 100.0, got %f", stats.ComplianceRate)
		}
	})
}

// The backend owns the status decision; the dashboard only renders it.
func TestReportSummaryClassify(t *testing.T) {
	cases := map[string]struct {
		total, compliant int
		wantStatus       string
		wantRate         string // "nil" or the formatted rate
	}{
		"null report": {0, 0, "empty", "nil"},
		"all pass":    {100, 100, "pass", "100.0"},
		"warn":        {100, 80, "warn", "80.0"},
		"fail":        {100, 79, "fail", "79.0"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			r := ReportSummary{TotalMessages: tc.total, CompliantMessages: tc.compliant}
			r.classify()

			gotRate := "nil"
			if r.ComplianceRate != nil {
				gotRate = fmt.Sprintf("%.1f", *r.ComplianceRate)
			}
			if r.Status != tc.wantStatus || gotRate != tc.wantRate {
				t.Fatalf("got status=%q rate=%s, want status=%q rate=%s", r.Status, gotRate, tc.wantStatus, tc.wantRate)
			}
		})
	}
}

// An RFC 7489 null report (one record, count=0, no auth results) must come
// back as status "empty" with a null compliance rate, not as a 0% failure.
func TestGetReports_NullReport(t *testing.T) {
	storage, err := NewStorage(":memory:")
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}
	defer func() { _ = storage.Close() }()

	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>example-reporter.tld</org_name>
    <email>dmarc-support@example-reporter.tld</email>
    <report_id>1780726747.321742392</report_id>
    <date_range>
      <begin>1780610400</begin>
      <end>1780696800</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>quarantine</p>
    <sp>quarantine</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip></source_ip>
      <count>0</count>
      <policy_evaluated>
        <disposition></disposition>
        <dkim></dkim>
        <spf></spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results></auth_results>
  </record>
</feedback>`

	feedback, err := parser.ParseReport([]byte(xmlData))
	if err != nil {
		t.Fatalf("Failed to parse report: %v", err)
	}
	if err := storage.SaveReport(feedback); err != nil {
		t.Fatalf("Failed to save report: %v", err)
	}

	reports, err := storage.GetReports(10, 0)
	if err != nil {
		t.Fatalf("Failed to get reports: %v", err)
	}
	if len(reports) != 1 {
		t.Fatalf("Expected 1 report, got %d", len(reports))
	}
	if reports[0].Status != "empty" || reports[0].ComplianceRate != nil {
		t.Fatalf("got status=%q rate=%v, want status=\"empty\" rate=nil", reports[0].Status, reports[0].ComplianceRate)
	}

	// The API contract the dashboard and MCP clients see.
	out, err := json.Marshal(reports[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"compliance_rate":null`, `"status":"empty"`} {
		if !strings.Contains(string(out), want) {
			t.Errorf("JSON %s\nmissing %s", out, want)
		}
	}
}
