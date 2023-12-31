package routing

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/redis/go-redis/v9"
)

type RoutingEngine struct {
	// Hold a reference to the Persistent Database
	db *sql.DB

	// Hold a reference to the Redis Routing Table
	rdb *redis.Client
}

func NewRoutingEngine(db *sql.DB, rdb *redis.Client) *RoutingEngine {
	return &RoutingEngine{
		db:  db,
		rdb: rdb,
	}
}

func (r *RoutingEngine) LaunchRoutingEngine() {
	// Prepopulate the routing table
	prepopulateRoutingTable(r.db, r.rdb)

	// Do processing here
	for {

		// TODO: Implement a way to check for expired records and remove them
		// from the routing table. Probably not MVP right now, dont have the
		// time.

		time.Sleep(1 * time.Second)
	}
}

func prepopulateRoutingTable(db *sql.DB, rdb *redis.Client) {

	// Get all records from the database
	records, err := database.GetDNSRecords(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Prepopulating Redis database with records, ", len(records))

	// Publish each record to the Redis database
	for _, record := range records {

		// Check if deactivated
		if !record.GetCommon().DeactivatedAt.IsZero() {
			continue
		}

		// Check if expired
		if !record.GetCommon().ExpiredAt.IsZero() && record.GetCommon().ExpiredAt.Before(time.Now()) {
			continue
		}

		fmt.Printf("Publishing [%s] %s record to Redis database\n", record.GetCommon().GetType(), record.GetCommon().GetFQDN())
		err = database.PublishRecordCache(rdb, record)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Callbacks for record changes, these will be used by the API to update the
// database

func (r *RoutingEngine) RecordCreateCb(rr models.RecordBlueprint) error {

	// Build the record
	record := rr.Build()

	// Insert the record into the database
	_, err := database.InsertDNSRecord(r.db, record)
	if err != nil {
		return err
	}

	return nil
}
func (r *RoutingEngine) RecordUpdateCb(rr models.RecordBlueprint) error {

	// Build the record
	record := rr.Build()

	// Update the record into the database
	return database.UpdateDNSRecord(r.db, record)
}
func (r *RoutingEngine) RecordDeleteCb(rr models.RecordBlueprint) error {

	if rr.Id == 0 {
		return fmt.Errorf("record does not exist")
	}

	// Delete the record from the database
	return database.DeleteDNSRecord(r.db, rr.Id)
}
