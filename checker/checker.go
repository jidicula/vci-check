package checker

type IssuerList struct {
	ParticipatingIssuers []Issuer `json:"participating_issuers"`
}

type Issuer struct {
	Iss          string `json:"iss"`
	Name         string `json:"name"`
	CanonicalIss string `json:"canonical_iss,omitempty"`
	Website      string `json:"website,omitempty"`
}
