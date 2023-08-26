package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/miekg/dns"
)

const (
	// DNS Record Types
	A     = int(dns.TypeA)
	CNAME = int(dns.TypeCNAME)
	SOA   = int(dns.TypeSOA)
	MX    = int(dns.TypeMX)
	TXT   = int(dns.TypeTXT)

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
	if r.Subdomain == "@" {
		return r.Domain + "."
	}
	return r.Subdomain + "." + r.Domain + "."
}

func (r RecordCommon) GetType() string {
	switch r.Type {
	case A:
		return "A"
	case CNAME:
		return "CNAME"
	case SOA:
		return "SOA"
	case MX:
		return "MX"
	case TXT:
		return "TXT"
	default:
		return "UNKNOWN"
	}
}

func ParseRecord(common RecordCommon, data []byte) (DNSRecord, error) {
	switch common.Type {
	case A:
		return unmarshalARecord(common, data)
	case CNAME:
		return unmarshalCNAMERecord(common, data)
	case MX:
		return unmarshalMXRecord(common, data)
	case TXT:
		return unmarshalTXTRecord(common, data)
	case SOA:
		return unmarshalSOARecord(common, data)
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
	case CNAME:
		var CNAMERecord CNAMERecord
		err = json.Unmarshal(data, &CNAMERecord)
		return &CNAMERecord, err
	case MX:
		var MXRecord MXRecord
		err = json.Unmarshal(data, &MXRecord)
		return &MXRecord, err
	case TXT:
		var TXTRecord TXTRecord
		err = json.Unmarshal(data, &TXTRecord)
		return &TXTRecord, err
	case SOA:
		var SOARecord SOARecord
		err = json.Unmarshal(data, &SOARecord)
		return &SOARecord, err
	}

	return nil, fmt.Errorf("unsupported record type")
}

func MarshalJSON(record DNSRecord) ([]byte, error) {
	return json.Marshal(record)
}
