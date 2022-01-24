package checker

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestIssuerUnmarshal(t *testing.T) {
	testIssuer1 := `{
      "iss": "https://myvaccinerecord.cdph.ca.gov/creds",
      "name": "State of California"
    }`
	tests := map[string]struct {
		data []byte
		want Issuer
	}{
		"issuer": {
			data: []byte(testIssuer1),
			want: Issuer{
				Iss:  "https://myvaccinerecord.cdph.ca.gov/creds",
				Name: "State of California",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := Issuer{}
			err := json.Unmarshal(tt.data, &got)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: got %v, want %v", name, got, tt.want)
			}
		})
	}

}

var testIssuerList = `{
  "participating_issuers": [
    {
      "iss": "https://myvaccinerecord.cdph.ca.gov/creds",
      "name": "State of California"
    },
    {
      "iss": "https://healthcardcert.lawallet.com",
      "name": "State of Louisiana"
    }
  ]
}`

func TestIssuerListUnmarshal(t *testing.T) {
	tests := map[string]struct {
		data []byte
		want IssuerList
	}{
		"2sameIssuer": {data: []byte(testIssuerList), want: IssuerList{
			ParticipatingIssuers: []Issuer{
				{
					Iss:  "https://myvaccinerecord.cdph.ca.gov/creds",
					Name: "State of California",
				},
				{
					Iss:  "https://healthcardcert.lawallet.com",
					Name: "State of Louisiana",
				},
			},
		}},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := IssuerList{}
			err := json.Unmarshal(tt.data, &got)
			if err != nil {
				t.Fatalf(err.Error())
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s: got %v, want %v", name, got, tt.want)
			}
		})
	}
}

func TestIsTrusted(t *testing.T) {
	tests := map[string]struct {
		url  string
		want bool
	}{
		"trusted issuer": {
			url:  "https://myvaccinerecord.cdph.ca.gov/creds",
			want: true,
		},
		"untrusted issuer": {
			url:  "https://mallory.me/creds",
			want: false,
		},
	}
	il := IssuerList{
		ParticipatingIssuers: []Issuer{
			{
				Iss:  "https://myvaccinerecord.cdph.ca.gov/creds",
				Name: "State of California",
			},
			{
				Iss:  "https://healthcardcert.lawallet.com",
				Name: "State of Louisiana",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := il.IsTrusted(tt.url)
			if got != tt.want {
				t.Errorf("%s: got %v, want %v", name, got, tt.want)
			}
		})
	}
}
