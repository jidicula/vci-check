package checker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type IssuerList struct {
	ParticipatingIssuers []Issuer `json:"participating_issuers"`
}

// isTrusted checks if a provided issuer URL is in the list of trusted issuers.
func (il IssuerList) IsTrusted(issURL string) bool {
	for _, iss := range il.ParticipatingIssuers {
		if issURL == iss.Iss {
			return true
		}
	}
	return false
}

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
		return il, fmt.Errorf("Status not OK, %d and %s", resp.StatusCode, body)
	}

	err = json.Unmarshal(body, &il)
	return il, err
}
