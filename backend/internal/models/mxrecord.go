package models

import (
	"encoding/json"
	"time"

	"github.com/miekg/dns"
)

type MXRecord struct {
	RecordCommon

	Data MXRecordData `json:"data"`
}

type MXRecordData struct {
	MailServer string `json:"mail_server"`
	Preference uint16 `json:"preference"`
}

func unmarshalMXRecord(common RecordCommon, data []byte) (*MXRecord, error) {
	// Parse the record data
	var recordData MXRecordData
	err := json.Unmarshal(data, &recordData)
	if err != nil {
		return nil, err
	}

	return &MXRecord{
		RecordCommon: common,
		Data:         recordData,
	}, nil
}

func marshalMXRecordData(record *MXRecord) []byte {
	data, _ := json.Marshal(record.Data)
	return data
}

func (r *MXRecord) GetCommon() RecordCommon {
	return r.RecordCommon
}

func (r *MXRecord) GetData() []byte {
	return marshalMXRecordData(r)
}

// GetResponse returns a pre-processed DNS response for an A Record
func (r *MXRecord) GetResponse() dns.RR {
	// Check if the record is deactivated
	if !r.DeactivatedAt.IsZero() && time.Since(r.DeactivatedAt) > 0 {
		return nil
	}

	// Check if the record is expired
	if !r.ExpiredAt.IsZero() && time.Since(r.ExpiredAt) > 0 {
		return nil
	}

	return &dns.MX{
		Hdr: dns.RR_Header{
			Name:   r.GetFQDN(),
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    uint32(r.TTL),
		},
		Mx:         r.Data.MailServer,
		Preference: r.Data.Preference,
	}
}
