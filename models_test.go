package domainavailability

import (
	"encoding/json"
	"testing"
)

// TestTime tests JSON encoding/parsing functions for the StringBool values
func TestStringBool(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"UNAVAILABLE"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"AVAILABLE"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "API error: [] \"\" is unexpected value for domainAvailability",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v *StringBool

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

// checkErr checks for an error.
func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
