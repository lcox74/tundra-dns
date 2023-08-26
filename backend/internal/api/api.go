package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lcox74/tundra-dns/backend/internal/database"
)

func LaunchRouter(db *sql.DB) {
	r := mux.NewRouter()

	// API Route GET /api/records
	r.HandleFunc("/api/records", func(w http.ResponseWriter, r *http.Request) {
		GetRecords(db, w, r)
	}).Methods("GET")

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
