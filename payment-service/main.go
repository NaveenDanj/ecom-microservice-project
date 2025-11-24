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
	db, err = sql.Open("postgres", "host=localhost port=5436 user=payment password=password dbname=payment_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/create-payment", RegisterHandler).Methods("POST")
	fmt.Println("Payment service running on port 8084")
	log.Fatal(http.ListenAndServe(":8084", router))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Payment added successfully!"))
}
