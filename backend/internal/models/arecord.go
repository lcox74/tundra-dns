package models

import (
	"net"
	"time"

	"github.com/miekg/dns"
)

type ARecord struct {
	RecordCommon

	// Target Resource
	Address string `json:"address"`

	// Fallback (Optional)
	LastSeen time.Time `json:"last_seen"`
	Fallback string    `json:"fallback,omitempty"`
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *ARecord) GetResponse() dns.RR {
	// Check if the target is deactivated
	if time.Since(r.LastSeen) > DefaultLastSeenTTL*time.Second {
		return nil
	}

	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	// Check if the record has been seen recently
	if r.Fallback != "" && time.Since(r.LastSeen) > DefaultLastSeenTTL*time.Second {
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   r.GetFQDN(),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			A: net.ParseIP(r.Fallback),
		}
	}

	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		A: net.ParseIP(r.Address),
	}
}
