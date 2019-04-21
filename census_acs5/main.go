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
	Geoid10 string			//make this dynamic 04.08.19?
	Fields int  					//scan  / numeric in golang ()   .scan function
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

	subject := req.PathParameters["subject"]
	geounit := req.PathParameters["geounit"]
	year := req.PathParameters["year"]
	// also need some form of strong params?
	// geoid := req.PathParameters["geoid10"] // request params
	// fields := req.PathParameters["fields"]
	// query := req.URL.Query()
	// // fmt.Println(req.URL.String())
	//e.g. api.shiftresearchlab.org/census/acs5/{subject}/{geounit}/{year}?geoid10='08001'&fields=b01001_001e,b01001_002e

//	q, err := url.ParseQuery(req.Body)	//
	// if err != nil {
	// 	return Response{StatusCode: 400, Headers: headers}, errors.Wrap(err, "Bad input")
	// }


	geoid10 := req.QueryStringParameters["geoid10"]
	fields := req.QueryStringParameters["fields"]

	//get count
	//var count int

	fmt.Println(geoid10)
	fmt.Println(fields)

	tableString := "acs5." + geounit + "_state_" + subject + "_" + year	//double check format for subject

	rows, err := db.Query("SELECT geoid10, " + fields + " FROM " + tableString " WHERE geoid10")
	//rows, err := db.Query("SELECT geoid10, b01001_001e FROM acs5.county_state_b01001_2016") //original

	defer rows.Close()
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	for rows.Next() {
		line := item{}
		err = rows.Scan(&line.Geoid10)
		err = rows.Scan(&line.Fields)
		// for i := 0; i < count; i++ {					// use reg for
		// 	var field int
		// 	err = rows.Scan(&field)
		// 	line.Fields = append(line.Fields, field)			///
		// }
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
