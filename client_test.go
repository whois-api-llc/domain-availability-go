package domainavailability

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathDomainAvailabilityResponseOK         = "/DomainAvailability/ok"
	pathDomainAvailabilityResponseError      = "/DomainAvailability/error"
	pathDomainAvailabilityResponse500        = "/DomainAvailability/500"
	pathDomainAvailabilityResponsePartial1   = "/DomainAvailability/partial"
	pathDomainAvailabilityResponsePartial2   = "/DomainAvailability/partial2"
	pathDomainAvailabilityResponseUnparsable = "/DomainAvailability/unparsable"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

// dummyServer is the sample of the Domain Availability API server for testing.
func dummyServer(resp, respUnparsable string, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathDomainAvailabilityResponseOK:
		case pathDomainAvailabilityResponseError:
			w.WriteHeader(499)
			response = respErr
		case pathDomainAvailabilityResponse500:
			w.WriteHeader(500)
			response = respUnparsable
		case pathDomainAvailabilityResponsePartial1:
			response = response[:len(response)-10]
		case pathDomainAvailabilityResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		case pathDomainAvailabilityResponseUnparsable:
			response = respUnparsable
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

// newAPI returns new Domain Availability API client for testing.
func newAPI(apiServer *httptest.Server, link string) *Client {
	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}

	apiURL.Path = link

	params := ClientParams{
		HTTPClient:                apiServer.Client(),
		DomainAvailabilityBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

// TestDomainAvailabilityGet tests the Get function.
func TestDomainAvailabilityGet(t *testing.T) {
	checkResultRec := func(res *DomainAvailabilityResponse) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"DomainInfo":{"domainAvailability":"UNAVAILABLE","domainName":"whoisxmlapi.com"}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"errorCode":"WHOIS_00","msg":"Test error message."}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDomainAvailabilityResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDomainAvailabilityResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "partial response 1",
			path: pathDomainAvailabilityResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathDomainAvailabilityResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathDomainAvailabilityResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "API error: [WHOIS_00] Test error message.",
		},
		{
			name: "unparsable response",
			path: pathDomainAvailabilityResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: "cannot parse response: invalid character '<' looking for beginning of value",
		},
		{
			name: "invalid argument",
			path: pathDomainAvailabilityResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			want:    false,
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			gotRec, _, err := api.Get(tt.args.ctx, tt.args.options.mandatory, tt.args.options.option)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DomainAvailability.Get() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.want {
				if !checkResultRec(gotRec) {
					t.Errorf("DomainAvailability.Get() got = %v, expected something else", gotRec)
				}
			} else {
				if gotRec != nil {
					t.Errorf("DomainAvailability.Get() got = %v, expected nil", gotRec)
				}
			}
		})
	}
}

// TestDomainAvailabilityGetRaw tests the GetRaw function.
func TestDomainAvailabilityGetRaw(t *testing.T) {
	checkResultRaw := func(res []byte) bool {
		return len(res) != 0
	}

	ctx := context.Background()

	const resp = `{"DomainInfo":{"domainAvailability":"UNAVAILABLE","domainName":"whoisxmlapi.com"}}`

	const respUnparsable = `<?xml version="1.0" encoding="utf-8"?><>`

	const errResp = `{"ErrorMessage":{"errorCode":"WHOIS_00","msg":"Test error message."}}`

	server := dummyServer(resp, respUnparsable, errResp)
	defer server.Close()

	type options struct {
		mandatory string
		option    Option
	}

	type args struct {
		ctx     context.Context
		options options
	}

	tests := []struct {
		name    string
		path    string
		args    args
		wantErr string
	}{
		{
			name: "successful request",
			path: pathDomainAvailabilityResponseOK,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathDomainAvailabilityResponse500,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 500",
		},
		{
			name: "partial response 1",
			path: pathDomainAvailabilityResponsePartial1,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "partial response 2",
			path: pathDomainAvailabilityResponsePartial2,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "unparsable response",
			path: pathDomainAvailabilityResponseUnparsable,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "",
		},
		{
			name: "could not process request",
			path: pathDomainAvailabilityResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"whoisxmlapi.com",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: "API failed with status code: 499",
		},
		{
			name: "invalid argument",
			path: pathDomainAvailabilityResponseError,
			args: args{
				ctx: ctx,
				options: options{
					"",
					OptionOutputFormat("JSON"),
				},
			},
			wantErr: `invalid argument: "domainName" can not be empty`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newAPI(server, tt.path)

			resp, err := api.GetRaw(tt.args.ctx, tt.args.options.mandatory)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("DomainAvailability.GetRaw() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if resp != nil && !checkResultRaw(resp.Body) {
				t.Errorf("DomainAvailability.GetRaw() got = %v, expected something else", string(resp.Body))
			}
		})
	}
}
