package imap

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/meysam81/parse-dmarc/internal/config"
)

// Client represents an IMAP client
type Client struct {
	config *config.IMAPConfig
	client *client.Client
}

// NewClient creates a new IMAP client
func NewClient(cfg *config.IMAPConfig) *Client {
	return &Client{config: cfg}
}

// Connect establishes connection to IMAP server
func (c *Client) Connect() error {
	var imapClient *client.Client
	var err error

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	log.Printf("Connecting to %s...", addr)

	if c.config.UseTLS {
		imapClient, err = client.DialTLS(addr, &tls.Config{
			ServerName: c.config.Host,
		})
	} else {
		imapClient, err = client.Dial(addr)
	}

	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.client = imapClient
	log.Printf("Connected to %s", addr)

	// Login
	if err := c.client.Login(c.config.Username, c.config.Password); err != nil {
		_ = c.client.Logout()
		return fmt.Errorf("login failed: %w", err)
	}

	log.Printf("Logged in as %s", c.config.Username)
	return nil
}

// Disconnect closes the IMAP connection
func (c *Client) Disconnect() error {
	if c.client != nil {
		return c.client.Logout()
	}
	return nil
}

// Report represents a DMARC report email
type Report struct {
	Subject     string
	From        string
	Date        string
	Attachments []Attachment
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Data     []byte
}

// FetchDMARCReports fetches DMARC reports from the mailbox
func (c *Client) FetchDMARCReports() ([]Report, error) {
	// Select mailbox
	mbox, err := c.client.Select(c.config.Mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select mailbox: %w", err)
	}

	if mbox.Messages == 0 {
		log.Println("No messages in mailbox")
		return []Report{}, nil
	}

	// Search for unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	ids, err := c.client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(ids) == 0 {
		log.Println("No new messages found")
		return []Report{}, nil
	}

	log.Printf("Found %d new messages", len(ids))

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchFlags}

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	reports := []Report{}

	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			log.Printf("Server didn't return message body for UID %d", msg.Uid)
			continue
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Printf("Failed to create mail reader: %v", err)
			continue
		}

		report := Report{
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date.String(),
		}

		if len(msg.Envelope.From) > 0 {
			report.From = msg.Envelope.From[0].Address()
		}

		// Process email parts
		for {
			part, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Error reading part: %v", err)
				break
			}

			switch h := part.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				// Only process DMARC-related attachments
				if isDMARCAttachment(filename) {
					data, err := io.ReadAll(part.Body)
					if err != nil {
						log.Printf("Error reading attachment: %v", err)
						continue
					}

					report.Attachments = append(report.Attachments, Attachment{
						Filename: filename,
						Data:     data,
					})
				}
			}
		}

		// Only add reports with attachments
		if len(report.Attachments) > 0 {
			reports = append(reports, report)
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	return reports, nil
}

// MarkAsSeen marks messages as seen
func (c *Client) MarkAsSeen(messageIDs []uint32) error {
	if len(messageIDs) == 0 {
		return nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(messageIDs...)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}

	return c.client.Store(seqSet, item, flags, nil)
}

// isDMARCAttachment checks if filename is likely a DMARC report
func isDMARCAttachment(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".xml") ||
		strings.HasSuffix(lower, ".xml.gz") ||
		strings.HasSuffix(lower, ".zip") ||
		strings.Contains(lower, "dmarc")
}
