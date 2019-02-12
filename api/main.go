package main

import (
	"encoding/json"
	"log"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "guest"
	DB_PASSWORD = "guest"
	DB_NAME     = "shift_public"
	DB_HOST     = "demodb.catj63cigq6x.us-east-2.rds.amazonaws.com"
)

func main() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	log.Println("successful connection to db")

	type item struct {
		Geoid10 string
		Field   int
	}

	type dat struct {
		Items []item
	}

	output := dat{}

	rows, err := db.Query("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016")
	defer rows.Close()
	if err != nil {
		log.Fatal(err.Error())
	}

	for rows.Next() {
		line := item{}
		err = rows.Scan(
			&line.Geoid10,
			&line.Field,
		)
		output.Items = append(output.Items, line)
	}

	fmt.Printf("%+v\n", output)

	_, err = json.Marshal(output)
	if err != nil {
		log.Fatal(err.Error())
	}
}
