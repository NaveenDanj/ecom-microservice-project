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
	db, err = sql.Open("postgres", "host=localhost port=5434 user=product password=password dbname=product_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	router.HandleFunc("/add", RegisterHandler).Methods("POST")

	fmt.Println("Product service running on port 8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Product added successfully!"))
}
