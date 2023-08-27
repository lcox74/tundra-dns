package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lcox74/tundra-dns/backend/internal/database"
	"github.com/lcox74/tundra-dns/backend/internal/models"
	"github.com/lcox74/tundra-dns/backend/internal/routing"
)

func LaunchRouter(db *sql.DB, engine *routing.RoutingEngine) {
	r := mux.NewRouter()

	// API Route GET /api/records
	r.HandleFunc("/api/records", func(w http.ResponseWriter, r *http.Request) {
		GetRecords(db, w, r)
	}).Methods("GET")

	// API Route POST /api/record
	r.HandleFunc("/api/record", func(w http.ResponseWriter, r *http.Request) {
		CreateRecord(engine.RecordCreateCb, w, r)
	}).Methods("POST")

	corsObj := handlers.AllowedOrigins([]string{"*"})
	http.ListenAndServe(":8053", handlers.CORS(corsObj)(r))
}

func GetRecords(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetRecords")

	rawRecords, err := database.GetDNSRecords(db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Convert to JSON
	jsonRecords, err := json.Marshal(rawRecords)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Write JSON to ResponseWriter
	w.Write(jsonRecords)
}

func CreateRecord(cb func(models.RecordBlueprint) error, w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateRecord")

	// Decode the JSON
	var recordBlueprint models.RecordBlueprint
	err := json.NewDecoder(r.Body).Decode(&recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	// Create the record
	err = cb(recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusCreated)
}
