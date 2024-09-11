package conf

type ACME struct {
	Name             string
	Mode             string            `json:"CertMode"` // file, http, dns
	RejectUnknownSni bool              `json:"RejectUnknownSni"`
	Provider         string            `json:"Provider"` // alidns, cloudflare, gandi, godaddy....
	Email            string            `json:"Email"`
	DNSEnv           map[string]string `json:"DNSEnv"`
	Storage          string            `json:"Storage"`
}
