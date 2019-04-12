package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"regexp"
)

const (
	DB_USER     = "guest"
	DB_PASSWORD = "guest"
	DB_NAME     = "shift_public"
	DB_HOST     = "demodb.catj63cigq6x.us-east-2.rds.amazonaws.com"
)

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

	//rows, err := db.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = 'geography'")
	//
	//for rows.Next() {
	//	tableName := ""
	//	rows.Scan(&tableName)
	//	fields := strings.Split(tableName, "_")
	//	fmt.Println(fields)
	//}
	//
	//defer rows.Close()

	var output interface{}

	geoUnit := req.PathParameters["geounit"]

	geoParam := ""

	year := req.PathParameters["year"]

	if matches, _ := regexp.MatchString("[\\d]{3}0", year); !matches {
		return Response{StatusCode: 422}, errors.New("invalid year: " + year)
	}
	if geoUnit == "nbhd" {
		geoParam = "nhid"
	} else if geoUnit == "county" || geoUnit == "tract" || geoUnit == "block"{
		geoParam = "geoid" + year[len(year) - 2:]
	} else {
		return Response{StatusCode: 422}, errors.New("invalid geographical unit: " + geoUnit)
	}

	geoParamValue := req.QueryStringParameters[geoParam]

	tableName := "geography." + geoUnit + "_state_geography_" + year

	query := fmt.Sprintf("SELECT json_build_object('geometry', ST_AsGeoJSON(geometry)::json, 'properties', json_build_object('" + geoParam + "', " + geoParam + ")) FROM %s WHERE %s = '%s'", tableName, geoParam, geoParamValue)

	rows, err := db.Query(query)

	defer rows.Close()
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	for rows.Next() {
		var scannedRow []byte
		err = rows.Scan(&scannedRow)
		if err != nil {
			return Response{StatusCode: 500}, err
		}
		err = json.Unmarshal(scannedRow, &output)
		if err != nil {
			return Response{StatusCode: 500}, err
		}
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
