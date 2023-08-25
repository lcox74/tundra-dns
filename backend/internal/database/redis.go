package database

import (
	"context"
	"fmt"

	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func InitialiseRedisDb(host string, port string) (*redis.Client, error) {
	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // No password set
		DB:       0,  // Default DB
	})

	if rdb == nil {
		return nil, fmt.Errorf("failed to initialise redis database")
	}

	return rdb, nil
}

func generateHashKey(recordType string, fqdn string) string {
	return fmt.Sprintf("%s:%s", recordType, fqdn)
}

func FetchRecordCache(rdb *redis.Client, recordType string, fqdn string) (models.DNSRecord, error) {
	// Generate the hash key
	hashKey := generateHashKey(recordType, fqdn)

	fmt.Println("Fetch:", hashKey)

	// Get the record from the database
	record, err := rdb.Get(ctx, hashKey).Result()
	if err != nil {
		return nil, err
	}

	// Unmarshal the record
	return models.UnmarshalJSON([]byte(record))
}

func PublishRecordCache(rdb *redis.Client, record models.DNSRecord) error {
	// Generate the hash key
	hashKey := generateHashKey(record.GetCommon().GetType(), record.GetCommon().GetFQDN())

	fmt.Println("Publish:", hashKey)

	// Marshal the record
	recordData, err := models.MarshalJSON(record)
	if err != nil {
		return err
	}

	// Publish the record to the database
	err = rdb.Set(ctx, hashKey, recordData, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
