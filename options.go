package domainavailability

import (
	"net/url"
	"strings"
)

// Option adds parameters to the query.
type Option func(v url.Values)

var _ = []Option{
	OptionOutputFormat("JSON"),
	OptionMode("DNS_ONLY"),
	OptionCredits("WHOIS"),
}

// OptionOutputFormat sets Response output format JSON | XML. Default: JSON.
func OptionOutputFormat(outputFormat string) Option {
	return func(v url.Values) {
		v.Set("outputFormat", strings.ToUpper(outputFormat))
	}
}

// OptionMode sets the check mode. The default mode is the fastest, the DNS_AND_WHOIS mode is slower but more accurate.
// Acceptable values: DNS_AND_WHOIS|DNS_ONLY. Default: DNS_ONLY.
func OptionMode(mode string) Option {
	return func(v url.Values) {
		v.Set("mode", strings.ToUpper(mode))
	}
}

// OptionCredits sets the type of credits used.
// DA — Domain Availability API credits will be taken into account when the API is called.
// WHOIS — WHOIS API credits will be taken into account when the API is called.
// Acceptable values: DA|WHOIS. Default: WHOIS.
func OptionCredits(credits string) Option {
	return func(v url.Values) {
		v.Set("credits", strings.ToUpper(credits))
	}
}
