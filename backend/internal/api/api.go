package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

	// API Route DELETE /api/record/{id}
	r.HandleFunc("/api/record/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL
		vars := mux.Vars(r)
		id := vars["id"]

		DeleteRecord(engine.RecordDeleteCb, w, r, id)
	}).Methods("DELETE")

	// API Route PUT /api/record/{id}
	r.HandleFunc("/api/record/{id}", func(w http.ResponseWriter, r *http.Request) {
		// Get the ID from the URL
		vars := mux.Vars(r)
		id := vars["id"]

		UpdateRecord(engine.RecordUpdateCb, db, w, r, id)
	}).Methods("PUT")

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
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	// Create the record
	err = cb(recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusCreated)
}

func DeleteRecord(cb func(models.RecordBlueprint) error, w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("DeleteRecord")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	// id to number
	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	var recordBlueprint models.RecordBlueprint
	recordBlueprint.Id = idNum
	err = cb(recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record Deleted"))
}

func UpdateRecord(cb func(models.RecordBlueprint) error, db *sql.DB, w http.ResponseWriter, r *http.Request, id string) {
	fmt.Println("UpdateRecord")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID"))
		return
	}

	// id to number
	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	// Get the record from the database
	rawRecord, err := database.GetDNSRecord(db, idNum)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// Decode the JSON
	var recordBlueprint models.RecordBlueprint
	err = json.NewDecoder(r.Body).Decode(&recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}
	recordBlueprint.Id = idNum
	recordBlueprint.Override(rawRecord)

	err = cb(recordBlueprint)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		fmt.Println(err)
		return
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Record Updated"))
}
