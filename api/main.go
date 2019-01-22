package main

import (
	"log"

	"fmt"
	"database/sql"
	_ "github.com/lib/pq"

)

const (
	DB_USER = "guest"
	DB_PASSWORD = "guest"
	DB_NAME = "shift_public"
	DB_HOST = "demodb.catj63cigq6x.us-east-2.rds.amazonaws.com"
)

func main() {
  connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	log.Println("successful connection to db")

	var geo string
	var b01001 int
	row := db.QueryRow("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016 WHERE geoid10='08001'")
	err = row.Scan(&geo, &b01001)
	if err != nil {
		log.Fatal(err.Error())
	}
	
	log.Println(b01001)
	log.Println(geo)

}