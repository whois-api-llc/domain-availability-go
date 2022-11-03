package domainavailability

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DomainAvailabilityService is an interface for Domain Availability API.
type DomainAvailabilityService interface {
	// Get returns parsed Domain Availability API response.
	Get(ctx context.Context, domainName string, opts ...Option) (*DomainAvailabilityResponse, *Response, error)

	// GetRaw returns raw Domain Availability API response as the Response struct with Body saved as a byte slice.
	GetRaw(ctx context.Context, domainName string, opts ...Option) (*Response, error)
}

// Response is the http.Response wrapper with Body saved as a byte slice.
type Response struct {
	*http.Response

	// Body is the byte slice representation of http.Response Body
	Body []byte
}

// domainAvailabilityServiceOp is the type implementing the DomainAvailability interface.
type domainAvailabilityServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ DomainAvailabilityService = &domainAvailabilityServiceOp{}

// newRequest creates the API request with default parameters and the specified apiKey.
func (service domainAvailabilityServiceOp) newRequest() (*http.Request, error) {
	req, err := service.client.NewRequest(http.MethodGet, service.baseURL, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

// apiResponse is used for parsing Domain Availability API response as a model instance.
type apiResponse struct {
	DomainAvailabilityResponse `json:"DomainInfo"`
	ErrorMessage               `json:"ErrorMessage"`
}

// request returns intermediate API response for further actions.
func (service domainAvailabilityServiceOp) request(ctx context.Context, domainName string, opts ...Option) (*Response, error) {
	if domainName == "" {
		return nil, &ArgError{"domainName", "can not be empty"}
	}

	req, err := service.newRequest()
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("domainName", domainName)

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b bytes.Buffer

	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return &Response{
			Response: resp,
			Body:     b.Bytes(),
		}, err
	}

	return &Response{
		Response: resp,
		Body:     b.Bytes(),
	}, nil
}

// parse parses raw Domain Availability API response.
func parse(raw []byte) (*apiResponse, error) {
	var response apiResponse

	err := json.NewDecoder(bytes.NewReader(raw)).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("cannot parse response: %w", err)
	}

	return &response, nil
}

// Get returns parsed Domain Availability API response.
func (service domainAvailabilityServiceOp) Get(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (domainAvailabilityResponse *DomainAvailabilityResponse, resp *Response, err error) {
	optsJSON := make([]Option, 0, len(opts)+1)
	optsJSON = append(optsJSON, opts...)
	optsJSON = append(optsJSON, OptionOutputFormat("JSON"))

	resp, err = service.request(ctx, domainName, optsJSON...)
	if err != nil {
		return nil, resp, err
	}

	domainAvailabilityResp, err := parse(resp.Body)
	if err != nil {
		return nil, resp, err
	}

	if domainAvailabilityResp.Message != "" || domainAvailabilityResp.Code != "" {
		return nil, nil, &ErrorMessage{
			Code:    domainAvailabilityResp.Code,
			Message: domainAvailabilityResp.Message,
		}
	}

	return &domainAvailabilityResp.DomainAvailabilityResponse, resp, nil
}

// GetRaw returns raw Domain Availability API response as the Response struct with Body saved as a byte slice.
func (service domainAvailabilityServiceOp) GetRaw(
	ctx context.Context,
	domainName string,
	opts ...Option,
) (resp *Response, err error) {
	resp, err = service.request(ctx, domainName, opts...)
	if err != nil {
		return resp, err
	}

	if respErr := checkResponse(resp.Response); respErr != nil {
		return resp, respErr
	}

	return resp, nil
}

// ArgError is the argument error.
type ArgError struct {
	Name    string
	Message string
}

// Error returns error message as a string.
func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
