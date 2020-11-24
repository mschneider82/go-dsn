package dsn

import (
	"errors"

	"golang.org/x/net/idna"
	"golang.org/x/text/unicode/norm"
)

var (
	ErrUnicodeMailbox = errors.New("address: cannot convert the Unicode local-part to the ACE form")
)

// toASCII converts the domain part of the email address to the A-label form and
// fails with ErrUnicodeMailbox if the local-part contains non-ASCII characters.
func toASCII(addr string) (string, error) {
	mbox, domain, err := split(addr)
	if err != nil {
		return addr, err
	}

	for _, ch := range mbox {
		if ch > 128 {
			return addr, ErrUnicodeMailbox
		}
	}

	if domain == "" {
		return mbox, nil
	}

	aDomain, err := idna.ToASCII(domain)
	if err != nil {
		return addr, err
	}

	return mbox + "@" + aDomain, nil
}

// toUnicode converts the domain part of the email address to the U-label form.
func toUnicode(addr string) (string, error) {
	mbox, domain, err := split(addr)
	if err != nil {
		return norm.NFC.String(addr), err
	}

	if domain == "" {
		return mbox, nil
	}

	uDomain, err := idna.ToUnicode(domain)
	if err != nil {
		return norm.NFC.String(addr), err
	}

	return mbox + "@" + norm.NFC.String(uDomain), nil
}

// addrSelectIDNA is a convenience function for conversion of domains in the email
// addresses to/from the Punycode form.
//
// ulabel=true => ToUnicode is used.
// ulabel=false => ToASCII is used.
func addrSelectIDNA(ulabel bool, addr string) (string, error) {
	if ulabel {
		return toUnicode(addr)
	}
	return toASCII(addr)
}

// dnsSelectIDNA is a convenience function for encoding to/from Punycode.
//
// If ulabel is true, it returns U-label encoded domain in the Unicode NFC
// form.
// If ulabel is false, it returns A-label encoded domain.
func dnsSelectIDNA(ulabel bool, domain string) (string, error) {
	if ulabel {
		uDomain, err := idna.ToUnicode(domain)
		return norm.NFC.String(uDomain), err
	}
	return idna.ToASCII(domain)
}
