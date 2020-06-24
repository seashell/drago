package domain

const (
	TokenTypeClient     = "client"
	TokenTypeManagement = "management"
)

type Token struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Subject   string `json:"subject"`
	Raw       string `json:"secret"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
	NotBefore int64  `json:"notBefore"`
}
