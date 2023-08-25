package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

const (
	// DNS Record Types
	A = int(dns.TypeA)

	// DNS Record Route Types
	Single   = 0
	Weighted = 1

	// DNS Record TTLs
	DefaultTTLSec      = 300
	DefaultLastSeenTTL = 60
)

type DNSRecord interface {
	GetCommon() RecordCommon
	GetData() []byte
	GetResponse() dns.RR
}

type RecordCommon struct {
	Id        int    `json:"id"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`

	// DNS Record Types
	Type      int `json:"type"`
	RouteType int `json:"route_type"`

	// Times Data
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeactivatedAt time.Time `json:"deactivated_at"`
	ExpiredAt     time.Time `json:"expired_at"`
	TTL           int       `json:"ttl"` // Time to Live (in seconds)

	// Allow and Deny Lists (IPs and Subnets)
	Allow []string `json:"allow"`
	Deny  []string `json:"deny"`
}

func (r RecordCommon) GetFQDN() string {
	return r.Subdomain + "." + r.Domain + "."
}

func (r RecordCommon) GetType() string {
	switch r.Type {
	case A:
		return "A"
	default:
		return "UNKNOWN"
	}
}

func ParseRecord(common RecordCommon, data []byte) (DNSRecord, error) {
	switch common.Type {
	case A:
		return unmarshalARecord(common, data)
	default:
		return nil, fmt.Errorf("unsupported record type")
	}
}

func UnmarshalJSON(data []byte) (DNSRecord, error) {
	var common RecordCommon
	err := json.Unmarshal(data, &common)
	if err != nil {
		return nil, err
	}

	switch common.Type {
	case A:
		var ARecord ARecord
		err = json.Unmarshal(data, &ARecord)
		return &ARecord, err
	}

	return nil, fmt.Errorf("unsupported record type")
}

func MarshalJSON(record DNSRecord) ([]byte, error) {
	return json.Marshal(record)
}
