package api

import (
	"encoding/json"

	"github.com/lcox74/tundra-dns/backend/internal/models"
)

type RecordBlueprint struct {
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	TTL       int    `json:"ttl"`

	Data interface{} `json:"data"`
}

func (r RecordBlueprint) Build() models.DNSRecord {

	// Create the Common
	common := models.RecordCommon{
		Id:        r.Id,
		Domain:    r.Domain,
		Subdomain: r.Subdomain,
		TTL:       r.TTL,
		RouteType: models.Single,
	}

	switch r.Type {
	case "A":
		d := models.ARecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.A
		return &models.ARecord{
			RecordCommon: common,
			Data:         d,
		}
	case "CNAME":
		d := models.CNAMERecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.CNAME
		return &models.CNAMERecord{
			RecordCommon: common,
			Data:         d,
		}
	case "SOA":
		d := models.SOARecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.SOA
		return &models.SOARecord{
			RecordCommon: common,
			Data:         d,
		}
	case "MX":
		d := models.MXRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.MX
		return &models.MXRecord{
			RecordCommon: common,
			Data:         d,
		}
	case "NS":
		d := models.NSRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.NS
		return &models.NSRecord{
			RecordCommon: common,
			Data:         d,
		}
	case "TXT":
		d := models.TXTRecordData{}
		jsonData, _ := json.Marshal(r.Data)
		json.Unmarshal(jsonData, &d)

		common.Type = models.TXT
		return &models.TXTRecord{
			RecordCommon: common,
			Data:         d,
		}
	default:
		return nil
	}
}
