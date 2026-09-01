package storage

import (
	"testing"

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

		if stats.EnforcedMessages != 0 {
			t.Errorf("Expected EnforcedMessages to be 0, got %d", stats.EnforcedMessages)
		}
	})

	t.Run("blocked spoofing counts as enforced", func(t *testing.T) {
		xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<feedback>
  <version>1.0</version>
  <report_metadata>
    <org_name>google.com</org_name>
    <email>noreply-dmarc-support@google.com</email>
    <report_id>12345678901234567891</report_id>
    <date_range>
      <begin>1609545600</begin>
      <end>1609632000</end>
    </date_range>
  </report_metadata>
  <policy_published>
    <domain>example.com</domain>
    <adkim>r</adkim>
    <aspf>r</aspf>
    <p>reject</p>
    <sp>reject</sp>
    <pct>100</pct>
  </policy_published>
  <record>
    <row>
      <source_ip>198.51.100.7</source_ip>
      <count>40</count>
      <policy_evaluated>
        <disposition>reject</disposition>
        <dkim>fail</dkim>
        <spf>fail</spf>
      </policy_evaluated>
    </row>
    <identifiers>
      <header_from>example.com</header_from>
    </identifiers>
    <auth_results>
      <spf>
        <domain>example.com</domain>
        <result>fail</result>
      </spf>
    </auth_results>
  </record>
  <record>
    <row>
      <source_ip>192.0.2.1</source_ip>
      <count>10</count>
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

		// 100 from the first subtest + 50 from this report
		if stats.TotalMessages != 150 {
			t.Errorf("Expected TotalMessages to be 150, got %d", stats.TotalMessages)
		}

		if stats.CompliantMessages != 110 {
			t.Errorf("Expected CompliantMessages to be 110, got %d", stats.CompliantMessages)
		}

		if stats.EnforcedMessages != 40 {
			t.Errorf("Expected EnforcedMessages to be 40, got %d", stats.EnforcedMessages)
		}
	})
}
