package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
)

const (
	DB_USER     = "guest"
	DB_PASSWORD = "guest"
	DB_NAME     = "shift_public"
	DB_HOST     = "demodb.catj63cigq6x.us-east-2.rds.amazonaws.com"
)

type item struct {
	Geoid10 string
	Field   int
}

type data struct {
	Items []item
}

type Response events.APIGatewayProxyResponse

func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s", DB_USER, DB_PASSWORD, DB_NAME, DB_HOST)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	defer db.Close()

	output := data{}

	subject := req.PathParameters["subject"] // need to know what subject is
	geounit := req.PathParameters["geounit"]
	year := req.PathParameters["year"]
	// also need some form of strong params
	// geoid := req.PathParameters["geoid10"] // request params
	// fields := req.PathParameters["fields"]
	// query := req.URL.Query()
	// fmt.Println(req.URL.String())

	q, err := url.ParseQuery(req.Body)
	if err != nil {
		return Response{StatusCode: 400, Headers: headers}, errors.Wrap(err, "Bad input")
	}

	geoid10 := q.Get("geoid10")
	fields := q.Get("fields")

	tableString := "acs5." + geounit + "_" + subject + "_" + year	//double check format for subject

	//rows, err := db.Query("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016") //original
	rows, err := db.Query("SELECT" +  geoid10 + ", " + fields + " FROM " + tableString)

	defer rows.Close()
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	for rows.Next() {
		line := item{}
		err = rows.Scan(
			&line.Geoid10,
			&line.Field,
		)
		output.Items = append(output.Items, line)
	}

	jsonB, err := json.Marshal(output)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// output is json data structure
	return Response{StatusCode: 200, Headers: headers, Body: string(jsonB)}, nil
}

func main() {
	lambda.Start(Handler)
}
