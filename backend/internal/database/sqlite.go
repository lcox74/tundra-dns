package database

import (
	"database/sql"
	"strings"
	"time"

	"github.com/lcox74/tundra-dns/backend/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

func InitialiseSqliteDb(dbPath string) (*sql.DB, error) {
	// Open the SQLite3 database.
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create the tables if they don't exist.
	err = initilaiseSqliteTables(db)

	return db, err
}

func initilaiseSqliteTables(db *sql.DB) error {

	// SQL statements to create tables if they don't exist.
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,

			domain TEXT NOT NULL,
			subdomain TEXT NOT NULL,
			ttl INTEGER NOT NULL,

			type INTEGER NOT NULL,
			route_type INTEGER NOT NULL,

			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deactivated_at DATETIME,
			expired_at DATETIME,

			allow BLOB,
			deny BLOB,
			data BLOB NOT NULL
		);
	`

	// Execute the SQL statements.
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func GetDNSRecord(db *sql.DB, id int) (models.DNSRecord, error) {

	// SQL query to select the record with the given ID.
	query := `
        SELECT * FROM records WHERE id = ?
    `

	// Query the database.
	row := db.QueryRow(query, id)
	return scanToRecord(row)
}
func GetDNSRecordFQDN(db *sql.DB, fqdn string) (models.DNSRecord, error) {

	// SQL query to select the record with the given ID.
	query := `
        SELECT * FROM records WHERE subdomain || '.' || domain = ?
    `

	// Query the database.
	row := db.QueryRow(query, fqdn)
	return scanToRecord(row)
}

func GetDNSRecords(db *sql.DB) ([]models.DNSRecord, error) {

	// SQL query to select all records.
	query := `
		SELECT * FROM records
	`

	// Query the database.
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan the rows into records.
	records := make([]models.DNSRecord, 0)
	for rows.Next() {
		rows.Scan()
		record, err := scanToRecordRows(rows)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

func InsertDNSRecord(db *sql.DB, record models.DNSRecord) (int, error) {
	// SQL query to insert a record.
	query := `
		INSERT INTO records (
			domain,
			subdomain,
			ttl,
			type,
			route_type,
			allow,
			deny,
			data
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Parse the record data.
	common := record.GetCommon()
	data := record.GetData()

	// Execute the query.
	result, err := db.Exec(
		query,
		common.Domain,
		common.Subdomain,
		common.TTL,
		common.Type,
		common.RouteType,
		[]byte(strings.Join(common.Allow, ",")),
		[]byte(strings.Join(common.Deny, ",")),
		data,
	)
	if err != nil {
		return 0, err
	}

	// Get the ID of the inserted record.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func scanToRecord(row *sql.Row) (models.DNSRecord, error) {
	var (
		createdAt     string
		updatedAt     string
		deactivatedAt sql.NullString
		expiredAt     sql.NullString
		allow         []byte
		deny          []byte
		data          []byte
	)

	recordCommon := models.RecordCommon{
		DeactivatedAt: time.Time{},
		ExpiredAt:     time.Time{},
	}

	// Scan the row into variables.
	err := row.Scan(
		&recordCommon.Id,
		&recordCommon.Domain,
		&recordCommon.Subdomain,
		&recordCommon.TTL,
		&recordCommon.Type,
		&recordCommon.RouteType,
		&createdAt,
		&updatedAt,
		&deactivatedAt,
		&expiredAt,
		&allow,
		&deny,
		&data,
	)
	if err != nil {
		return nil, err
	}

	// Parse the timestamps.
	recordCommon.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	recordCommon.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	// Parse the nullable timestamps.
	if deactivatedAt.Valid {
		recordCommon.DeactivatedAt, _ = time.Parse(time.RFC3339, deactivatedAt.String)
	}
	if expiredAt.Valid {
		recordCommon.ExpiredAt, _ = time.Parse(time.RFC3339, expiredAt.String)
	}

	// Parse the allow and deny lists.
	recordCommon.Allow = strings.Split(string(allow), ",")
	recordCommon.Deny = strings.Split(string(deny), ",")

	return models.ParseRecord(recordCommon, data)
}

func scanToRecordRows(rows *sql.Rows) (models.DNSRecord, error) {
	var (
		createdAt     string
		updatedAt     string
		deactivatedAt sql.NullString
		expiredAt     sql.NullString
		allow         []byte
		deny          []byte
		data          []byte
	)

	recordCommon := models.RecordCommon{
		DeactivatedAt: time.Time{},
		ExpiredAt:     time.Time{},
	}

	// Scan the row into variables.
	err := rows.Scan(
		&recordCommon.Id,
		&recordCommon.Domain,
		&recordCommon.Subdomain,
		&recordCommon.TTL,
		&recordCommon.Type,
		&recordCommon.RouteType,
		&createdAt,
		&updatedAt,
		&deactivatedAt,
		&expiredAt,
		&allow,
		&deny,
		&data,
	)
	if err != nil {
		return nil, err
	}

	// Parse the timestamps.
	recordCommon.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	recordCommon.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	// Parse the nullable timestamps.
	if deactivatedAt.Valid {
		recordCommon.DeactivatedAt, _ = time.Parse(time.RFC3339, deactivatedAt.String)
	}
	if expiredAt.Valid {
		recordCommon.ExpiredAt, _ = time.Parse(time.RFC3339, expiredAt.String)
	}

	// Parse the allow and deny lists.
	recordCommon.Allow = strings.Split(string(allow), ",")
	recordCommon.Deny = strings.Split(string(deny), ",")

	return models.ParseRecord(recordCommon, data)
}
