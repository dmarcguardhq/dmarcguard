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
	// A raw (uncompressed) report gets the same ceiling as a decompressed one.
	if len(data) > maxDecompressedSize {
		return nil, errReportTooLarge
	}

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

// tryDecompress attempts to decompress data (gzip or zip). Detection is
// decisive so that a valid-but-oversized member is reported as an error
// instead of silently falling through to "return the raw bytes" (which would
// hand a bomb to the XML parser and re-buffer it).
func tryDecompress(data []byte) ([]byte, error) {
	// gzip mandates its magic bytes, so this check is exact.
	if bytes.HasPrefix(data, []byte{0x1f, 0x8b}) {
		return decompressGzip(data)
	}

	// Let archive/zip decide: it locates the central directory from the end,
	// so archives with a leading stub (self-extracting) are recognised too,
	// which a fixed "PK\x03\x04" prefix check would miss.
	if zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data))); err == nil {
		return decompressZip(zipReader)
	}

	// Not compressed.
	return data, nil
}

// maxDecompressedSize bounds a single report, raw or decompressed. DMARC
// aggregate reports are small (usually well under a megabyte); a generous
// 16 MB ceiling stops a malicious sender from turning a tiny gzip/zip member
// into a multi-gigabyte allocation (a decompression bomb).
const maxDecompressedSize = 16 * 1024 * 1024

var errReportTooLarge = fmt.Errorf("report exceeds %d bytes", maxDecompressedSize)

// readAllLimited reads from r but rejects streams that exceed
// maxDecompressedSize instead of allocating without bound. The buffer starts
// small so normal reports don't pay for the cap; growth stops at the cap, so a
// bomb costs at most a few times maxDecompressedSize in transient allocation.
func readAllLimited(r io.Reader) ([]byte, error) {
	// Read one byte past the cap so we can tell "exactly at the cap" from
	// "over the cap".
	buf := bytes.NewBuffer(make([]byte, 0, 64*1024))
	n, err := io.Copy(buf, io.LimitReader(r, maxDecompressedSize+1))
	if err != nil {
		return nil, err
	}
	if n > maxDecompressedSize {
		return nil, errReportTooLarge
	}
	return buf.Bytes(), nil
}

// decompressGzip decompresses gzip data
func decompressGzip(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()

	return readAllLimited(reader)
}

// decompressZip returns the first XML member of the archive.
func decompressZip(zipReader *zip.Reader) ([]byte, error) {
	// Pick the first XML member; archives may carry directories or __MACOSX noise.
	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(file.Name), ".xml") {
			continue
		}
		// Reject a member whose declared uncompressed size is already over the
		// cap before opening it. archive/zip stops a read that exceeds the
		// declared size with ErrFormat, so an understated header can't bypass this.
		if file.UncompressedSize64 > maxDecompressedSize {
			return nil, errReportTooLarge
		}
		rc, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer func() { _ = rc.Close() }()

		return readAllLimited(rc)
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
