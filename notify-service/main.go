package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "host=localhost port=5437 user=notify password=password dbname=notification_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/create-notification", RegisterHandler).Methods("GET")

	fmt.Println("Notify service running on port 8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Notification created successfully!"))
}
