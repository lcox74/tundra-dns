package models

import (
	"encoding/json"
	"time"

	"github.com/miekg/dns"
)

type CNAMERecord struct {
	RecordCommon

	Data CNAMERecordData `json:"data"`
}

type CNAMERecordData struct {
	Alias string `json:"alias"`
}

func unmarshalCNAMERecord(common RecordCommon, data []byte) (*CNAMERecord, error) {
	// Parse the record data
	var recordData CNAMERecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &CNAMERecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalCNAMERecordData(record *CNAMERecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *CNAMERecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *CNAMERecord) GetData() []byte {
	return marshalCNAMERecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *CNAMERecord) GetResponse() dns.RR {
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
		Target: r.Data.Alias,
	}
}
