package imap

import (
	"crypto/tls"
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"github.com/meysam81/parse-dmarc/internal/config"
	"github.com/rs/zerolog"
)

// Client represents an IMAP client
type Client struct {
	config *config.IMAPConfig
	client *client.Client
	log    *zerolog.Logger
}

// NewClient creates a new IMAP client
func NewClient(cfg *config.IMAPConfig, log *zerolog.Logger) *Client {
	return &Client{config: cfg, log: log}
}

// Connect establishes connection to IMAP server
func (c *Client) Connect() error {
	var imapClient *client.Client
	var err error

	addr := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	c.log.Debug().Str("addr", addr).Msg("connecting")

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
	c.log.Info().Str("addr", addr).Msg("connected")

	// Login
	if err := c.client.Login(c.config.Username, c.config.Password); err != nil {
		_ = c.client.Logout()
		return fmt.Errorf("login failed: %w", err)
	}

	c.log.Info().Str("username", c.config.Username).Msg("logged in")
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

// FetchResult holds the fetched reports and the message sequence numbers
type FetchResult struct {
	Reports    []Report
	MessageIDs []uint32
}

// FetchDMARCReports fetches DMARC reports from the mailbox
func (c *Client) FetchDMARCReports() (*FetchResult, error) {
	// Select mailbox
	mbox, err := c.client.Select(c.config.Mailbox, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select mailbox: %w", err)
	}

	if mbox.Messages == 0 {
		c.log.Info().Msg("no messages in mailbox")
		return &FetchResult{}, nil
	}

	// Search for unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}

	ids, err := c.client.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	if len(ids) == 0 {
		c.log.Info().Msg("no new messages found")
		return &FetchResult{}, nil
	}

	c.log.Info().Int("count", len(ids)).Msg("found new messages")

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	section := &imap.BodySectionName{Peek: !c.config.MarkAsSeen}
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchFlags}

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	var reports []Report
	var messageIDs []uint32

	for msg := range messages {
		messageIDs = append(messageIDs, msg.SeqNum)

		r := msg.GetBody(section)
		if r == nil {
			c.log.Warn().Uint32("uid", msg.Uid).Msg("server didn't return message body")
			continue
		}

		mr, err := mail.CreateReader(r)
		if err != nil {
			c.log.Warn().Err(err).Msg("failed to create mail reader")
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
				c.log.Warn().Err(err).Msg("error reading part")
				break
			}

			switch h := part.Header.(type) {
			case *mail.AttachmentHeader:
				filename, _ := h.Filename()
				// Only process DMARC-related attachments
				if isDMARCAttachment(filename) {
					data, err := io.ReadAll(part.Body)
					if err != nil {
						c.log.Warn().Err(err).Msg("error reading attachment")
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

	return &FetchResult{Reports: reports, MessageIDs: messageIDs}, nil
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

// MoveMessages moves messages to a destination mailbox using COPY + DELETE + EXPUNGE
func (c *Client) MoveMessages(messageIDs []uint32, destMailbox string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	// Ensure destination mailbox exists (ignore error if it already exists)
	if err := c.client.Create(destMailbox); err != nil {
		c.log.Debug().Err(err).Str("mailbox", destMailbox).Msg("create mailbox (may already exist)")
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(messageIDs...)

	// Copy to destination mailbox
	if err := c.client.Copy(seqSet, destMailbox); err != nil {
		return fmt.Errorf("copy to %s: %w", destMailbox, err)
	}

	// Mark originals as deleted
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := c.client.Store(seqSet, item, flags, nil); err != nil {
		return fmt.Errorf("mark deleted: %w", err)
	}

	// Expunge deleted messages
	if err := c.client.Expunge(nil); err != nil {
		return fmt.Errorf("expunge: %w", err)
	}

	return nil
}

// isDMARCAttachment checks if filename is likely a DMARC report
func isDMARCAttachment(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".xml") ||
		strings.HasSuffix(lower, ".xml.gz") ||
		strings.HasSuffix(lower, ".zip") ||
		strings.Contains(lower, "dmarc")
}
