package models

import (
	"encoding/json"
	"net"
	"time"

	"github.com/miekg/dns"
)

type ARecord struct {
	RecordCommon

	Data ARecordData `json:"data"`
}

type ARecordData struct {
	Address string `json:"address"`

	// Fallback (Optional)
	LastSeen time.Time `json:"last_seen,omitempty"`
	Fallback string    `json:"fallback,omitempty"`
}

func unmarshalARecord(common RecordCommon, data []byte) (*ARecord, error) {
	// Parse the record data
	var recordData ARecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &ARecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalARecordData(record *ARecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *ARecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *ARecord) GetData() []byte {
	return marshalARecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *ARecord) GetResponse() dns.RR {
	// // Check if the target is deactivated
	// if time.Since(r.Data.LastSeen) > DefaultLastSeenTTL*time.Second {
	// 	return nil
	// }

	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	// Check if the record has been seen recently
	if r.Data.Fallback != "" && time.Since(r.Data.LastSeen) > DefaultLastSeenTTL*time.Second {
		return &dns.A{
			Hdr: dns.RR_Header{
				Name:   r.GetFQDN(),
				Rrtype: dns.TypeA,
				Class:  dns.ClassINET,
				Ttl:    uint32(r.TTL),
			},
			A: net.ParseIP(r.Data.Fallback),
		}
	}

	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		A: net.ParseIP(r.Data.Address),
	}
}
