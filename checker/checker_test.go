// Copyright (C) 2022  Johanan Idicula
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

func TestNewIssuerList(t *testing.T) {
	t.Run("NewIssuerList", func(t *testing.T) {
		_, err := NewIssuerList()
		if err != nil {
			t.Error(err)
		}
	})
}
