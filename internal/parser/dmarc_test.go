package parser

import (
	"testing"
)

func TestParseReport(t *testing.T) {
	// Sample DMARC aggregate report XML (simplified)
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

	feedback, err := ParseReport([]byte(xmlData))
	if err != nil {
		t.Fatalf("Failed to parse report: %v", err)
	}

	// Verify metadata
	if feedback.ReportMetadata.OrgName != "google.com" {
		t.Errorf("Expected org_name google.com, got %s", feedback.ReportMetadata.OrgName)
	}

	if feedback.ReportMetadata.ReportID != "12345678901234567890" {
		t.Errorf("Expected report_id 12345678901234567890, got %s", feedback.ReportMetadata.ReportID)
	}

	// Verify policy
	if feedback.PolicyPublished.Domain != "example.com" {
		t.Errorf("Expected domain example.com, got %s", feedback.PolicyPublished.Domain)
	}

	if feedback.PolicyPublished.P != "none" {
		t.Errorf("Expected policy none, got %s", feedback.PolicyPublished.P)
	}

	// Verify records
	if len(feedback.Records) != 1 {
		t.Fatalf("Expected 1 record, got %d", len(feedback.Records))
	}

	record := feedback.Records[0]
	if record.Row.SourceIP != "192.0.2.1" {
		t.Errorf("Expected source IP 192.0.2.1, got %s", record.Row.SourceIP)
	}

	if record.Row.Count != 100 {
		t.Errorf("Expected count 100, got %d", record.Row.Count)
	}

	if record.Row.PolicyEvaluated.DKIM != "pass" {
		t.Errorf("Expected DKIM pass, got %s", record.Row.PolicyEvaluated.DKIM)
	}

	if record.Row.PolicyEvaluated.SPF != "pass" {
		t.Errorf("Expected SPF pass, got %s", record.Row.PolicyEvaluated.SPF)
	}

	// Test helper methods
	totalMessages := feedback.GetTotalMessages()
	if totalMessages != 100 {
		t.Errorf("Expected total messages 100, got %d", totalMessages)
	}

	compliantCount := feedback.GetDMARCCompliantCount()
	if compliantCount != 100 {
		t.Errorf("Expected compliant count 100, got %d", compliantCount)
	}
}
