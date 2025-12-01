package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/goccy/go-json"

	"github.com/meysam81/parse-dmarc/internal/parser"
)

type Storage struct {
	db *sql.DB
}

type ReportSummary struct {
	ID                int64   `json:"id"`
	ReportID          string  `json:"report_id"`
	OrgName           string  `json:"org_name"`
	Domain            string  `json:"domain"`
	DateBegin         int64   `json:"date_begin"`
	DateEnd           int64   `json:"date_end"`
	TotalMessages     int     `json:"total_messages"`
	CompliantMessages int     `json:"compliant_messages"`
	ComplianceRate    float64 `json:"compliance_rate"`
	PolicyP           string  `json:"policy_p"`
}

type Statistics struct {
	TotalReports      int     `json:"total_reports"`
	TotalMessages     int     `json:"total_messages"`
	CompliantMessages int     `json:"compliant_messages"`
	ComplianceRate    float64 `json:"compliance_rate"`
	UniqueSourceIPs   int     `json:"unique_source_ips"`
	UniqueDomains     int     `json:"unique_domains"`
}

type TopSourceIP struct {
	SourceIP string `json:"source_ip"`
	Count    int    `json:"count"`
	Pass     int    `json:"pass"`
	Fail     int    `json:"fail"`
}

func (s *Storage) SaveReport(feedback *parser.Feedback) error {
	rawReport, err := json.Marshal(feedback)
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	result, err := tx.Exec(`
		INSERT OR IGNORE INTO reports (
			report_id, org_name, email, domain,
			date_begin, date_end, created_at,
			policy_p, policy_sp, policy_pct,
			total_messages, compliant_messages,
			raw_report
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		feedback.ReportMetadata.ReportID,
		feedback.ReportMetadata.OrgName,
		feedback.ReportMetadata.Email,
		feedback.PolicyPublished.Domain,
		feedback.ReportMetadata.DateRange.Begin,
		feedback.ReportMetadata.DateRange.End,
		time.Now().Unix(),
		feedback.PolicyPublished.P,
		feedback.PolicyPublished.SP,
		feedback.PolicyPublished.PCT,
		feedback.GetTotalMessages(),
		feedback.GetDMARCCompliantCount(),
		rawReport,
	)

	if err != nil {
		return fmt.Errorf("failed to insert report: %w", err)
	}

	reportID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil
	}

	for _, record := range feedback.Records {
		dkimDomains, _ := json.Marshal(record.AuthResults.DKIM)
		spfDomains, _ := json.Marshal(record.AuthResults.SPF)

		_, err := tx.Exec(`
			INSERT INTO records (
				report_id, source_ip, count,
				disposition, dkim_result, spf_result,
				header_from, envelope_from,
				dkim_domains, spf_domains
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			reportID,
			record.Row.SourceIP,
			record.Row.Count,
			record.Row.PolicyEvaluated.Disposition,
			record.Row.PolicyEvaluated.DKIM,
			record.Row.PolicyEvaluated.SPF,
			record.Identifiers.HeaderFrom,
			record.Identifiers.EnvelopeFrom,
			dkimDomains,
			spfDomains,
		)

		if err != nil {
			return fmt.Errorf("failed to insert record: %w", err)
		}
	}

	return tx.Commit()
}

func (s *Storage) GetReports(limit, offset int) ([]ReportSummary, error) {
	rows, err := s.db.Query(`
		SELECT id, report_id, org_name, domain,
		       date_begin, date_end,
		       total_messages, compliant_messages,
		       policy_p
		FROM reports
		ORDER BY date_begin DESC
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var reports []ReportSummary
	for rows.Next() {
		var r ReportSummary
		err := rows.Scan(
			&r.ID, &r.ReportID, &r.OrgName, &r.Domain,
			&r.DateBegin, &r.DateEnd,
			&r.TotalMessages, &r.CompliantMessages,
			&r.PolicyP,
		)
		if err != nil {
			return nil, err
		}

		if r.TotalMessages > 0 {
			r.ComplianceRate = float64(r.CompliantMessages) / float64(r.TotalMessages) * 100
		}

		reports = append(reports, r)
	}

	return reports, nil
}

func (s *Storage) GetReportByID(id int64) (*parser.Feedback, error) {
	var rawReport string
	err := s.db.QueryRow("SELECT raw_report FROM reports WHERE id = ?", id).Scan(&rawReport)
	if err != nil {
		return nil, err
	}

	var feedback parser.Feedback
	if err := json.Unmarshal([]byte(rawReport), &feedback); err != nil {
		return nil, err
	}

	return &feedback, nil
}

func (s *Storage) GetStatistics() (*Statistics, error) {
	var stats Statistics

	err := s.db.QueryRow(`
		SELECT
			COUNT(*) as total_reports,
			COALESCE(SUM(total_messages), 0) as total_messages,
			COALESCE(SUM(compliant_messages), 0) as compliant_messages
		FROM reports
	`).Scan(&stats.TotalReports, &stats.TotalMessages, &stats.CompliantMessages)

	if err != nil {
		return nil, err
	}

	if stats.TotalMessages > 0 {
		stats.ComplianceRate = float64(stats.CompliantMessages) / float64(stats.TotalMessages) * 100
	}

	err = s.db.QueryRow("SELECT COUNT(DISTINCT source_ip) FROM records").Scan(&stats.UniqueSourceIPs)
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow("SELECT COUNT(DISTINCT domain) FROM reports").Scan(&stats.UniqueDomains)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *Storage) GetTopSourceIPs(limit int) ([]TopSourceIP, error) {
	rows, err := s.db.Query(`
		SELECT
			source_ip,
			SUM(count) as total_count,
			SUM(CASE WHEN (dkim_result = 'pass' OR spf_result = 'pass') THEN count ELSE 0 END) as pass_count,
			SUM(CASE WHEN (dkim_result != 'pass' AND spf_result != 'pass') THEN count ELSE 0 END) as fail_count
		FROM records
		GROUP BY source_ip
		ORDER BY total_count DESC
		LIMIT ?
	`, limit)

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []TopSourceIP
	for rows.Next() {
		var r TopSourceIP
		if err := rows.Scan(&r.SourceIP, &r.Count, &r.Pass, &r.Fail); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

// DomainStats holds statistics for a single domain
type DomainStats struct {
	Domain            string  `json:"domain"`
	TotalMessages     int     `json:"total_messages"`
	CompliantMessages int     `json:"compliant_messages"`
	ComplianceRate    float64 `json:"compliance_rate"`
}

// OrgStats holds statistics for a reporting organization
type OrgStats struct {
	OrgName string `json:"org_name"`
	Reports int    `json:"reports"`
}

// DispositionStats holds statistics for a disposition type
type DispositionStats struct {
	Disposition string `json:"disposition"`
	Count       int    `json:"count"`
}

// AuthResultStats holds authentication result statistics
type AuthResultStats struct {
	Result string `json:"result"`
	Count  int    `json:"count"`
}

// GetDomainStats returns statistics grouped by domain
func (s *Storage) GetDomainStats() ([]DomainStats, error) {
	rows, err := s.db.Query(`
		SELECT domain,
		       COALESCE(SUM(total_messages), 0) as total_messages,
		       COALESCE(SUM(compliant_messages), 0) as compliant_messages
		FROM reports
		GROUP BY domain
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []DomainStats
	for rows.Next() {
		var ds DomainStats
		if err := rows.Scan(&ds.Domain, &ds.TotalMessages, &ds.CompliantMessages); err != nil {
			return nil, err
		}
		if ds.TotalMessages > 0 {
			ds.ComplianceRate = float64(ds.CompliantMessages) / float64(ds.TotalMessages) * 100
		}
		stats = append(stats, ds)
	}
	return stats, nil
}

// GetOrgStats returns statistics grouped by reporting organization
func (s *Storage) GetOrgStats() ([]OrgStats, error) {
	rows, err := s.db.Query(`
		SELECT org_name, COUNT(*) as reports
		FROM reports
		GROUP BY org_name
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []OrgStats
	for rows.Next() {
		var os OrgStats
		if err := rows.Scan(&os.OrgName, &os.Reports); err != nil {
			return nil, err
		}
		stats = append(stats, os)
	}
	return stats, nil
}

// GetDispositionStats returns message counts grouped by disposition
func (s *Storage) GetDispositionStats() ([]DispositionStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(disposition, 'unknown') as disposition,
		       SUM(count) as total_count
		FROM records
		GROUP BY disposition
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []DispositionStats
	for rows.Next() {
		var ds DispositionStats
		if err := rows.Scan(&ds.Disposition, &ds.Count); err != nil {
			return nil, err
		}
		stats = append(stats, ds)
	}
	return stats, nil
}

// GetSPFStats returns SPF authentication result statistics
func (s *Storage) GetSPFStats() ([]AuthResultStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(spf_result, 'unknown') as result,
		       SUM(count) as total_count
		FROM records
		GROUP BY spf_result
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []AuthResultStats
	for rows.Next() {
		var as AuthResultStats
		if err := rows.Scan(&as.Result, &as.Count); err != nil {
			return nil, err
		}
		stats = append(stats, as)
	}
	return stats, nil
}

// GetDKIMStats returns DKIM authentication result statistics
func (s *Storage) GetDKIMStats() ([]AuthResultStats, error) {
	rows, err := s.db.Query(`
		SELECT COALESCE(dkim_result, 'unknown') as result,
		       SUM(count) as total_count
		FROM records
		GROUP BY dkim_result
	`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var stats []AuthResultStats
	for rows.Next() {
		var as AuthResultStats
		if err := rows.Scan(&as.Result, &as.Count); err != nil {
			return nil, err
		}
		stats = append(stats, as)
	}
	return stats, nil
}
