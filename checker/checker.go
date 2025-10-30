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

// Package checker provides VCI types and functions for working with them.
package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// IssuerList is a struct containing a slice of trusted Issuers.
type IssuerList struct {
	ParticipatingIssuers []Issuer `json:"participating_issuers"`
}

// IsTrusted checks if a provided issuer URL is in the list of trusted issuers.
func (il IssuerList) IsTrusted(issURL string) bool {
	for _, iss := range il.ParticipatingIssuers {
		if issURL == iss.Iss {
			return true
		}
	}
	return false
}

// An Issuer is a struct corresponding to a vaccine credential issuer.
type Issuer struct {
	Iss          string `json:"iss"`
	Name         string `json:"name"`
	CanonicalIss string `json:"canonical_iss,omitempty"`
	Website      string `json:"website,omitempty"`
}

// NewIssuerList fetches the VCI Issuers JSON and unmarshals it into a new
// IssuerList.
func NewIssuerList() (IssuerList, error) {
	il := IssuerList{}

	resp, err := http.Get("https://raw.githubusercontent.com/the-commons-project/vci-directory/main/vci-issuers.json")
	if err != nil {
		return il, err
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return il, err
	}
	if resp.StatusCode != http.StatusOK {
		return il, fmt.Errorf("status not OK, %d and %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &il)
	return il, err
}
