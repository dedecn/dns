// Package dnsutil contains higher-level methods useful with the dns
// package.  While package dns implements the DNS protocols itself,
// these functions are related but not directly required for protocol
// processing.  They are often useful in preparing input/output of the
// functions in package dns.
package dnsutil

import (
	"strings"

	"github.com/dedecn/dns"
)

// AddOrigin adds origin to s if s is not already a FQDN.
// Note that the result may not be a FQDN.  If origin does not end
// with a ".", the result won't either.
// This implements the zonefile convention (specified in RFC 1035,
// Section "5.1. Format") that "@" represents the
// apex (bare) domain. i.e. AddOrigin("@", "foo.com.") returns "foo.com.".
func AddOrigin(s, origin string) string {
	// ("foo.", "origin.") -> "foo." (already a FQDN)
	// ("foo", "origin.") -> "foo.origin."
	// ("foo", "origin") -> "foo.origin"
	// ("foo", ".") -> "foo." (Same as dns.Fqdn())
	// ("foo.", ".") -> "foo." (Same as dns.Fqdn())
	// ("@", "origin.") -> "origin." (@ represents the apex (bare) domain)
	// ("", "origin.") -> "origin." (not obvious)
	// ("foo", "") -> "foo" (not obvious)

	if dns.IsFqdn(s) {
		return s // s is already a FQDN, no need to mess with it.
	}
	if origin == "" {
		return s // Nothing to append.
	}
	if s == "@" || s == "" {
		return origin // Expand apex.
	}
	if origin == "." {
		return dns.Fqdn(s)
	}

	return s + "." + origin // The simple case.
}

// TrimDomainName trims origin from s if s is a subdomain.
// This function will never return "", but returns "@" instead (@ represents the apex domain).
func TrimDomainName(s, origin string) string {
	// An apex (bare) domain is always returned as "@".
	// If the return value ends in a ".", the domain was not the suffix.
	// origin can end in "." or not. Either way the results should be the same.

	if s == "" {
		return "@"
	}
	// Someone is using TrimDomainName(s, ".") to remove a dot if it exists.
	if origin == "." {
		return strings.TrimSuffix(s, origin)
	}

	original := s
	s = dns.Fqdn(s)
	origin = dns.Fqdn(origin)

	if !dns.IsSubDomain(origin, s) {
		return original
	}

	slabels := dns.Split(s)
	olabels := dns.Split(origin)
	m := dns.CompareDomainName(s, origin)
	if len(olabels) == m {
		if len(olabels) == len(slabels) {
			return "@" // origin == s
		}
		if (s[0] == '.') && (len(slabels) == (len(olabels) + 1)) {
			return "@" // TrimDomainName(".foo.", "foo.")
		}
	}

	// Return the first (len-m) labels:
	return s[:slabels[len(slabels)-m]-1]
}
