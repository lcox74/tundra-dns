package main

import (
	"fmt"

	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/lcox74/tundra-dns/backend/internal/routing"
)

const TestFQDNSub = "test"
const TestFQDNDomain = "tundra.test"

func main() {
	fmt.Println("Hello, World!")

	// Initialise the database
	db, err := database.InitialiseSqliteDb("tundra-test.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	fqdn := fmt.Sprintf("%s.%s", TestFQDNSub, TestFQDNDomain)

	// Get record with FQDN "test.tundra.test"
	_, err = database.GetDNSRecordFQDN(db, fqdn)
	if err != nil {
		// Create A Record
		record := &models.ARecord{
			RecordCommon: models.RecordCommon{
				Domain:    TestFQDNDomain,
				Subdomain: TestFQDNSub,
				Type:      models.A,
				RouteType: models.Single,
				TTL:       models.DefaultTTLSec,
			},
			Address: "10.10.10.10",
		}

		// Insert the record into the database
		id, err := database.InsertDNSRecord(db, record)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Inserted record with ID: %d\n", id)
	}

	// Get the record from the database
	record, err := database.GetDNSRecordFQDN(db, fqdn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Got record: \t %s\t%s\n", record.GetCommon().GetType(), record.GetCommon().GetFQDN())

	// Launch the DNS Query Handler
	routing.LaunchDNSQueryHandler()
}
