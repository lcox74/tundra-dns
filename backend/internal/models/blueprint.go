package models

import (
	"encoding/json"
)

type RecordBlueprint struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	TTL       int    `json:"ttl"`

	Data interface{} `json:"data"`
}

func (r *RecordBlueprint) Override(rb DNSRecord) {
	if r.Domain == "" {
		r.Domain = rb.GetCommon().Domain
	}
	if r.Subdomain == "" {
		r.Subdomain = rb.GetCommon().Subdomain
	}
	if r.TTL == 0 {
		r.TTL = rb.GetCommon().TTL
	}
	if r.Type == "" {
		r.Type = rb.GetCommon().GetType()
	}
}

func (r RecordBlueprint) Build() DNSRecord {

	// Create the Common
	common := RecordCommon{
		Id:        r.Id,
		Domain:    r.Domain,
		Subdomain: r.Subdomain,
		TTL:       r.TTL,
		RouteType: Single,
	}

	switch r.Type {
	case "A":
		d := ARecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = A
		return &ARecord{
			RecordCommon: common,
			Data:         d,
		}
	case "CNAME":
		d := CNAMERecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = CNAME
		return &CNAMERecord{
			RecordCommon: common,
			Data:         d,
		}
	case "SOA":
		d := SOARecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = SOA
		return &SOARecord{
			RecordCommon: common,
			Data:         d,
		}
	case "MX":
		d := MXRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = MX
		return &MXRecord{
			RecordCommon: common,
			Data:         d,
		}
	case "NS":
		d := NSRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = NS
		return &NSRecord{
			RecordCommon: common,
			Data:         d,
		}
	case "TXT":
		d := TXTRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = TXT
		return &TXTRecord{
			RecordCommon: common,
			Data:         d,
		}
	default:
		return nil
	}
}
