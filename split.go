package dsn

import (
	"errors"
	"strings"
)

// split splits a email address (as defined by RFC 5321 as a forward-path
// token) into local part (mailbox) and domain.
//
// Note that definition of the forward-path token includes the special
// postmaster address without the domain part. split will return domain == ""
// in this case.
//
// split does almost no sanity checks on the input and is intentionally naive.
// If this is a concern, ValidMailbox and ValidDomain should be used on the
// output.
func split(addr string) (mailbox, domain string, err error) {
	if strings.EqualFold(addr, "postmaster") {
		return addr, "", nil
	}

	indx := strings.LastIndexByte(addr, '@')
	if indx == -1 {
		return "", "", errors.New("address: missing at-sign")
	}
	mailbox = addr[:indx]
	domain = addr[indx+1:]
	if mailbox == "" {
		return "", "", errors.New("address: empty local-part")
	}
	if domain == "" {
		return "", "", errors.New("address: empty domain")
	}
	return
}
