package domainavailability

import (
	"encoding/json"
	"fmt"
)

// unmarshalString parses the JSON-encoded data and returns value as a string.
func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// StringBool is a helper wrapper on bool
type StringBool bool

// UnmarshalJSON decodes AVAILABLE/UNAVAILABLE values from Domain Availability API.
func (b *StringBool) UnmarshalJSON(bytes []byte) error {
	str, err := unmarshalString(bytes)
	if err != nil {
		return err
	}

	switch str {
	case "AVAILABLE":
		*b = true
	case "UNAVAILABLE":
		*b = false
	default:
		return &ErrorMessage{"", `"` + str + `"` + " is unexpected value for domainAvailability"}
	}
	return nil
}

// MarshalJSON encodes AVAILABLE/UNAVAILABLE values to the Domain Availability API representation.
func (b *StringBool) MarshalJSON() ([]byte, error) {
	if *b {
		return []byte(`"AVAILABLE"`), nil
	}
	return []byte(`"UNAVAILABLE"`), nil
}

// DomainAvailabilityResponse is a response of Domain Availability API.
type DomainAvailabilityResponse struct {
	// DomainName is the target domain name.
	DomainName string `json:"domainName"`

	// IsAvailable is the registration state of the domain name.
	IsAvailable *StringBool `json:"domainAvailability"`
}

// ErrorMessage is the error message.
type ErrorMessage struct {
	Code    string `json:"errorCode"`
	Message string `json:"msg"`
}

// Error returns error message as a string.
func (e *ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%s] %s", e.Code, e.Message)
}
