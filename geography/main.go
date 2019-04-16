package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lib/pq"
	"regexp"
	"strings"
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
	
	geoUnit := req.PathParameters["geounit"]

	geoParam := ""

	year := req.PathParameters["year"]

	if matches, _ := regexp.MatchString("[\\d]{3}0", year); !matches {
		return Response{StatusCode: 422}, errors.New("invalid year: " + year)
	}
	if geoUnit == "nbhd" {
		geoParam = "nhid"
	} else if geoUnit == "county" || geoUnit == "tract" || geoUnit == "block" {
		geoParam = "geoid" + year[len(year)-2:]
	} else {
		return Response{StatusCode: 422}, errors.New("invalid geographical unit: " + geoUnit)
	}

	geoParamValues := strings.Split(req.QueryStringParameters[geoParam], ",")

	tableName := "geography." + geoUnit + "_state_geography_" + year

	//Query for geometries
	query := fmt.Sprintf("SELECT %[1]s, json_build_object('geography', ST_AsGeoJSON(ST_SIMPLIFYPRESERVETOPOLOGY(ST_COLLECT(geometry), .0001))::json) FROM %[2]s WHERE %[1]s = ANY($1) GROUP BY %[1]s", geoParam, tableName)

	rows, err := db.Query(query, pq.Array(geoParamValues))
	defer rows.Close()

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	//Generate FeatureCollection from GeometryCollection
	returnJSON := gabs.New()

	returnJSON.Set("FeatureCollection", "type")
	returnJSON.Array("features")

	for rows.Next() {
		var scannedID string
		var scannedGeometry []byte
		err = rows.Scan(&scannedID, &scannedGeometry)
		if err != nil {
			return Response{StatusCode: 500}, err
		}

		parsedJSON, err := gabs.ParseJSON(scannedGeometry)

		children, _ := parsedJSON.Path("geography.geometries").Children()

		for _, child := range children {
			featureJSON := gabs.New()
			featureJSON.Set("Feature", "type")
			featureJSON.Set(child.Data(), "geometry")
			featureJSON.SetP(scannedID, "properties."+geoParam)
			returnJSON.ArrayAppend(featureJSON.Data(), "features")
		}

		if err != nil {
			return Response{StatusCode: 500}, err
		}
	}

	jsonB, err := json.Marshal(returnJSON.Data())

	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// output is json data structure
	return Response{StatusCode: 200, Headers: headers, Body: string(jsonB)}, nil
}

func main() {
	lambda.Start(Handler)
}
