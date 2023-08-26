package models

import (
	"encoding/json"
	"time"

	"github.com/miekg/dns"
)

type SOARecord struct {
	RecordCommon

	Data SOARecordData `json:"data"`
}

type SOARecordData struct {
	Ns      string `json:"ns_server"`  // ns.tundra-dns.io
	Mbox    string `json:"ns_mailbox"` // admin.tundra-dns.io
	Serial  uint32 `json:"serial"`
	Refresh uint32 `json:"refresh"`
	Retry   uint32 `json:"retry"`
	Expire  uint32 `json:"expire"`
	Minttl  uint32 `json:"minttl"`
}

func unmarshalSOARecord(common RecordCommon, data []byte) (*SOARecord, error) {
	// Parse the record data
	var recordData SOARecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &SOARecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalSOARecordData(record *SOARecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *SOARecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *SOARecord) GetData() []byte {
	return marshalSOARecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *SOARecord) GetResponse() dns.RR {
	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	return &dns.SOA{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeSOA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Ns:      r.Data.Ns,
		Mbox:    r.Data.Mbox,
		Serial:  r.Data.Serial,
		Refresh: r.Data.Refresh,
		Retry:   r.Data.Retry,
		Expire:  r.Data.Expire,
		Minttl:  r.Data.Minttl,
	}
}
