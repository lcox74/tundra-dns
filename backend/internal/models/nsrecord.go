package models

import (
	"encoding/json"
	"time"

	"github.com/miekg/dns"
)

type NSRecord struct {
	RecordCommon

	Data NSRecordData `json:"data"`
}

type NSRecordData struct {
	NameServer string `json:"server"`
}

func unmarshalNSRecord(common RecordCommon, data []byte) (*NSRecord, error) {
	// Parse the record data
	var recordData NSRecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &NSRecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalNSRecordData(record *NSRecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *NSRecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *NSRecord) GetData() []byte {
	return marshalNSRecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *NSRecord) GetResponse() dns.RR {
	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Target: r.Data.NameServer,
	}
}
