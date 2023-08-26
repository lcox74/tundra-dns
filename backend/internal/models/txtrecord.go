package models

import (
	"encoding/json"
	"time"

	"github.com/miekg/dns"
)

type TXTRecord struct {
	RecordCommon

	Data TXTRecordData `json:"data"`
}

type TXTRecordData struct {
	Content []string `json:"content"`
}

func unmarshalTXTRecord(common RecordCommon, data []byte) (*TXTRecord, error) {
	// Parse the record data
	var recordData TXTRecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &TXTRecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalTXTRecordData(record *TXTRecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *TXTRecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *TXTRecord) GetData() []byte {
	return marshalTXTRecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *TXTRecord) GetResponse() dns.RR {
	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	return &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Txt: r.Data.Content,
	}
}
