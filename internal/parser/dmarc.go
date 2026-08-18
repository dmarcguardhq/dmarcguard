package parser

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"time"
)

// Feedback represents the root DMARC aggregate report structure (RFC 7489)
type Feedback struct {
	XMLName         xml.Name        `xml:"feedback"`
	Version         string          `xml:"version"`
	ReportMetadata  ReportMetadata  `xml:"report_metadata"`
	PolicyPublished PolicyPublished `xml:"policy_published"`
	Records         []Record        `xml:"record"`
}

// ReportMetadata contains information about the report
type ReportMetadata struct {
	OrgName          string    `xml:"org_name"`
	Email            string    `xml:"email"`
	ExtraContactInfo string    `xml:"extra_contact_info,omitempty"`
	ReportID         string    `xml:"report_id"`
	DateRange        DateRange `xml:"date_range"`
	Errors           []string  `xml:"error,omitempty"`
}

// DateRange specifies the time range for the report
type DateRange struct {
	Begin int64 `xml:"begin"`
	End   int64 `xml:"end"`
}

// PolicyPublished represents the DMARC policy as published in DNS
type PolicyPublished struct {
	Domain string `xml:"domain"`
	ADKIM  string `xml:"adkim,omitempty"` // DKIM alignment mode (r=relaxed, s=strict)
	ASPF   string `xml:"aspf,omitempty"`  // SPF alignment mode (r=relaxed, s=strict)
	P      string `xml:"p"`               // Policy (none, quarantine, reject)
	SP     string `xml:"sp,omitempty"`    // Subdomain policy
	PCT    int    `xml:"pct,omitempty"`   // Percentage of messages to filter
	FO     string `xml:"fo,omitempty"`    // Failure reporting options
}

// Record represents a single record in the aggregate report
type Record struct {
	Row         Row         `xml:"row"`
	Identifiers Identifiers `xml:"identifiers"`
	AuthResults AuthResults `xml:"auth_results"`
}

// Row contains policy evaluation results
type Row struct {
	SourceIP        string          `xml:"source_ip"`
	Count           int             `xml:"count"`
	PolicyEvaluated PolicyEvaluated `xml:"policy_evaluated"`
}

// PolicyEvaluated shows the result of policy evaluation
type PolicyEvaluated struct {
	Disposition string   `xml:"disposition"` // none, quarantine, reject
	DKIM        string   `xml:"dkim"`        // pass, fail
	SPF         string   `xml:"spf"`         // pass, fail
	Reason      []Reason `xml:"reason,omitempty"`
}

// Reason explains policy override
type Reason struct {
	Type    string `xml:"type"`
	Comment string `xml:"comment,omitempty"`
}

// Identifiers contains message identifiers
type Identifiers struct {
	EnvelopeTo   string `xml:"envelope_to,omitempty"`
	EnvelopeFrom string `xml:"envelope_from,omitempty"`
	HeaderFrom   string `xml:"header_from"`
}

// AuthResults contains authentication results
type AuthResults struct {
	DKIM []DKIMResult `xml:"dkim,omitempty"`
	SPF  []SPFResult  `xml:"spf"`
}

// DKIMResult represents DKIM authentication result
type DKIMResult struct {
	Domain      string `xml:"domain"`
	Selector    string `xml:"selector,omitempty"`
	Result      string `xml:"result"` // none, pass, fail, policy, neutral, temperror, permerror
	HumanResult string `xml:"human_result,omitempty"`
}

// SPFResult represents SPF authentication result
type SPFResult struct {
	Domain string `xml:"domain"`
	Scope  string `xml:"scope,omitempty"` // helo, mfrom
	Result string `xml:"result"`          // none, neutral, pass, fail, softfail, temperror, permerror
}

// ParseReport parses a DMARC aggregate report from raw data
func ParseReport(data []byte) (*Feedback, error) {
	// Try to decompress if needed
	decompressed, err := tryDecompress(data)
	if err != nil {
		return nil, fmt.Errorf("decompression failed: %w", err)
	}

	var feedback Feedback
	if err := xml.Unmarshal(decompressed, &feedback); err != nil {
		return nil, fmt.Errorf("XML parsing failed: %w", err)
	}

	return &feedback, nil
}

// tryDecompress attempts to decompress data (gzip or zip)
func tryDecompress(data []byte) ([]byte, error) {
	// Try gzip first
	if gzipData, err := decompressGzip(data); err == nil {
		return gzipData, nil
	}

	// Try zip
	if zipData, err := decompressZip(data); err == nil {
		return zipData, nil
	}

	// Return original data if not compressed
	return data, nil
}

// decompressGzip decompresses gzip data
func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	return io.ReadAll(reader)
}

// decompressZip decompresses zip data (returns first file)
func decompressZip(data []byte) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	// Pick the first XML member; archives may carry directories or __MACOSX noise.
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(file.Name), ".xml") {
			continue
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer func() { _ = rc.Close() }()

		return io.ReadAll(rc)
	}

	return nil, fmt.Errorf("no XML file in zip archive")
}

// GetDateRange returns the date range as time.Time objects
func (f *Feedback) GetDateRange() (time.Time, time.Time) {
	begin := time.Unix(f.ReportMetadata.DateRange.Begin, 0)
	end := time.Unix(f.ReportMetadata.DateRange.End, 0)
	return begin, end
}

// GetTotalMessages returns the total count of messages in the report
func (f *Feedback) GetTotalMessages() int {
	total := 0
	for _, record := range f.Records {
		total += record.Row.Count
	}
	return total
}

// GetDMARCCompliantCount returns count of DMARC-compliant messages
func (f *Feedback) GetDMARCCompliantCount() int {
	count := 0
	for _, record := range f.Records {
		if record.Row.PolicyEvaluated.DKIM == "pass" || record.Row.PolicyEvaluated.SPF == "pass" {
			count += record.Row.Count
		}
	}
	return count
}

// NormalizeForJSON ensures all slice fields are initialized (not nil) to produce
// valid JSON that matches the MCP output schema. The MCP SDK infers JSON schemas
// from Go types, and nil slices serialize as null which violates the array type
// requirement in the schema.
func (f *Feedback) NormalizeForJSON() {
	if f == nil {
		return
	}

	// Normalize ReportMetadata.Errors
	if f.ReportMetadata.Errors == nil {
		f.ReportMetadata.Errors = []string{}
	}

	// Normalize Records slice
	if f.Records == nil {
		f.Records = []Record{}
	}

	// Normalize each record's nested slices
	for i := range f.Records {
		// Normalize PolicyEvaluated.Reason
		if f.Records[i].Row.PolicyEvaluated.Reason == nil {
			f.Records[i].Row.PolicyEvaluated.Reason = []Reason{}
		}

		// Normalize AuthResults.DKIM
		if f.Records[i].AuthResults.DKIM == nil {
			f.Records[i].AuthResults.DKIM = []DKIMResult{}
		}

		// Normalize AuthResults.SPF
		if f.Records[i].AuthResults.SPF == nil {
			f.Records[i].AuthResults.SPF = []SPFResult{}
		}
	}
}
