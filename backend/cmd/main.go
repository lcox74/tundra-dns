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

	// Construct the FQDN
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
			Data: models.ARecordData{
				Address: "10.10.10.10",
			},
		}

		// Insert the record into the database
		id, err := database.InsertDNSRecord(db, record)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Inserted record with ID: %d\n", id)
	}
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
