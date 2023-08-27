package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lcox74/tundra-dns/backend/internal/api"
	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/lcox74/tundra-dns/backend/internal/routing"
)

const TestFQDNDomain = "tundra-test.xyz."

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

	// Create Routing Engine Here
	routingEngine := routing.NewRoutingEngine(db, rdb)

	// Launch the DNS Query Handler and API Router
	go api.LaunchRouter(db, routingEngine)
	go routing.LaunchDNSQueryHandler(rdb)

	// Launch Routing Engine
	routingEngine.LaunchRoutingEngine()
}

func populateRecords(db *sql.DB) {
	var err error

	// Check if there is a './initial.json' file
	_, err = os.Stat("/app/initial.json")
	if err != nil {
		fmt.Println("No initial.json file found.")
		return
	}

	// Load the initial.json file
	file, err := os.Open("/app/initial.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Parse the JSON file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse the JSON file
	var records []models.RecordBlueprint
	err = json.Unmarshal(data, &records)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Insert the records into the database
	for _, record := range records {
		// Insert the record into the database
		fmt.Println("Inserting: ", record)
		rr := record.Build()
		if rr == nil {
			fmt.Println("Failed to build record")
			continue
		}

		id, err := database.InsertDNSRecord(db, rr)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Inserted [%s] %s record with ID: %d\n", rr.GetCommon().GetType(), rr.GetCommon().GetFQDN(), id)
	}
}
