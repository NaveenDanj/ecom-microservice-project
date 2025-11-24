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
	db, err = sql.Open("postgres", "host=localhost port=5435 user=order password=password dbname=order_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/create-order", RegisterHandler).Methods("POST")

	fmt.Println("Order service running on port 8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Order added successfully!"))
}
