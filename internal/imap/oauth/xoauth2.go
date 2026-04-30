package oauth

import (
	"errors"

	"github.com/emersion/go-sasl"
)

// XOAUTH2 is the SASL mechanism name used by Gmail and Microsoft 365 for
// OAuth2-bearer IMAP/SMTP authentication. It is distinct from the IETF
// OAUTHBEARER mechanism (RFC 7628) — Microsoft 365 only advertises XOAUTH2
// for IMAP/SMTP, and Gmail prefers it.
const XOAUTH2 = "XOAUTH2"

// xoauth2Client implements sasl.Client for the XOAUTH2 mechanism.
// The wire format of the initial response is:
//
//	user={email}\x01auth=Bearer {accessToken}\x01\x01
type xoauth2Client struct {
	username string
	token    string
}

// NewXOAUTH2Client returns a SASL client for XOAUTH2. username is the email
// address the access token was issued to (XOAUTH2 binds the two — a mismatch
// causes Gmail/M365 to reject the connection).
func NewXOAUTH2Client(username, accessToken string) sasl.Client {
	return &xoauth2Client{username: username, token: accessToken}
}

func (c *xoauth2Client) Start() (mech string, ir []byte, err error) {
	ir = []byte("user=" + c.username + "\x01auth=Bearer " + c.token + "\x01\x01")
	return XOAUTH2, ir, nil
}

func (c *xoauth2Client) Next(challenge []byte) ([]byte, error) {
	// Server challenge on failure is a JSON error blob. Per RFC, the client
	// must send an empty response to acknowledge the failure; the server
	// then sends a tagged BAD/NO. We surface the challenge bytes in the
	// error so callers/logs see the provider's diagnostic.
	if len(challenge) > 0 {
		return []byte(""), errors.New("xoauth2 server challenge: " + string(challenge))
	}
	return nil, errors.New("xoauth2: unexpected empty server challenge")
}
