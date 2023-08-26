package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lcox74/tundra-dns/backend/internal/api"
	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/lcox74/tundra-dns/backend/internal/routing"
	"github.com/redis/go-redis/v9"
)

const TestFQDNDomain = "tundra-test.xyz"

func main() {
	fmt.Println("Hello, World!")

	// Initialise the database
	db, err := database.InitialiseSqliteDb("tundra-test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Create a new Redis client
	rdb, err := database.InitialiseRedisDb(os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rdb.Close()

	// Lets set some dummy data
	populateRecords(db)
	populateRoutingTable(db, rdb)

	// TODO: Create Routing Engine Here...

	// Launch the DNS Query Handler and API Router
	go api.LaunchRouter(db)
	routing.LaunchDNSQueryHandler(rdb)

	// TODO: Launch Routing Engine Here...
}

func populateRecords(db *sql.DB) {
	var err error

	// Get record with FQDN "test.tundra.test"

	// Create A Record
	record1 := &models.ARecord{
		RecordCommon: models.RecordCommon{
			Domain:    TestFQDNDomain,
			Subdomain: "@",
			Type:      models.A,
			RouteType: models.Single,
			TTL:       models.DefaultTTLSec,
		},
		Data: models.ARecordData{
			Address: "10.10.10.10",
		},
	}

	// Insert the record into the database
	id, err := database.InsertDNSRecord(db, record1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted A record with ID: %d\n", id)

	// Create SOA Record
	record2 := &models.SOARecord{
		RecordCommon: models.RecordCommon{
			Domain:    TestFQDNDomain,
			Subdomain: "@",
			Type:      models.SOA,
			RouteType: models.Single,
			TTL:       models.DefaultTTLSec,
		},
		Data: models.SOARecordData{
			Ns:      "ns.tundra-dns.io.",
			Mbox:    "admin.tundra-dns.io.",
			Serial:  2318336624,
			Refresh: 10000,
			Retry:   2400,
			Expire:  604800,
			Minttl:  1800,
		},
	}

	// Insert the record into the database
	id, err = database.InsertDNSRecord(db, record2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted SOA record with ID: %d\n", id)

	// Create CNAME Record
	record3 := &models.CNAMERecord{
		RecordCommon: models.RecordCommon{
			Domain:    TestFQDNDomain,
			Subdomain: "www",
			Type:      models.CNAME,
			RouteType: models.Single,
			TTL:       models.DefaultTTLSec,
		},
		Data: models.CNAMERecordData{
			Alias: TestFQDNDomain,
		},
	}

	// Insert the record into the database
	id, err = database.InsertDNSRecord(db, record3)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted CNAME record with ID: %d\n", id)

	// Create MX Record
	record4 := &models.MXRecord{
		RecordCommon: models.RecordCommon{
			Domain:    TestFQDNDomain,
			Subdomain: "@",
			Type:      models.MX,
			RouteType: models.Single,
			TTL:       models.DefaultTTLSec,
		},
		Data: models.MXRecordData{
			MailServer: "mail.tundra-dns.io",
			Preference: 10,
		},
	}

	// Insert the record into the database
	id, err = database.InsertDNSRecord(db, record4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted MX record with ID: %d\n", id)

	// Create TXT Record
	record5 := &models.TXTRecord{
		RecordCommon: models.RecordCommon{
			Domain:    TestFQDNDomain,
			Subdomain: "@",
			Type:      models.TXT,
			RouteType: models.Single,
			TTL:       models.DefaultTTLSec,
		},
		Data: models.TXTRecordData{
			Content: []string{"Hello", "World"},
		},
	}

	// Insert the record into the database
	id, err = database.InsertDNSRecord(db, record5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Inserted TXT record with ID: %d\n", id)

}

func populateRoutingTable(db *sql.DB, rdb *redis.Client) {

	// Get all records from the database
	records, err := database.GetDNSRecords(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Publish each record to the Redis database
	for _, record := range records {
		fmt.Println("Publishing record to Redis database ", record.GetCommon().GetFQDN())
		err = database.PublishRecordCache(rdb, record)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
