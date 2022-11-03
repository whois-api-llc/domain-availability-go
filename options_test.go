package domainavailability

import (
	"net/url"
	"reflect"
	"testing"
)

// TestOptions tests the Options functions.
func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "outputFormat1",
			values: url.Values{},
			option: OptionOutputFormat("JSON"),
			want:   "outputFormat=JSON",
		},
		{
			name:   "outputFormat2",
			values: url.Values{},
			option: OptionOutputFormat("xml"),
			want:   "outputFormat=XML",
		},
		{
			name:   "mode1",
			values: url.Values{},
			option: OptionMode("DNS_ONLY"),
			want:   "mode=DNS_ONLY",
		},
		{
			name:   "mode2",
			values: url.Values{},
			option: OptionMode("dns_and_whois"),
			want:   "mode=DNS_AND_WHOIS",
		},
		{
			name:   "credits1",
			values: url.Values{},
			option: OptionCredits("DA"),
			want:   "credits=DA",
		},
		{
			name:   "credits2",
			values: url.Values{},
			option: OptionCredits("whois"),
			want:   "credits=WHOIS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
