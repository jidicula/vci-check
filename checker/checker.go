package checker

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
