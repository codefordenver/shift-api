package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/Jeffail/gabs"
	"regexp"
	"strconv"
	"strings"
)

const (
	DB_USER     = "guest"
	DB_PASSWORD = "guest"
	DB_NAME     = "shift_public"
	DB_HOST     = "demodb.catj63cigq6x.us-east-2.rds.amazonaws.com"
)

type Response events.APIGatewayProxyResponse

type Entry struct {
	Geoid  string
	Fields []Field
}

type Field struct {
	Field string
	Value int
}

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

	subject := req.PathParameters["subject"]

	if matches, _ := regexp.MatchString("b[\\w]{5}", subject); !matches {
		return Response{StatusCode: 422}, errors.New("invalid subject: " + subject)
	}

	fieldPattern := subject + "_[\\d]{3}[a-z]"

	geoUnit := req.PathParameters["geounit"]

	if !(geoUnit == "tract" || geoUnit == "county") {
		return Response{StatusCode: 422}, errors.New("invalid geounit: " + geoUnit)
	}

	year := req.PathParameters["year"]

	if matches, _ := regexp.MatchString("[\\d]{4}", year); !matches {
		return Response{StatusCode: 422}, errors.New("invalid year: " + year)
	}

	geoParam := "geoid" + string([]byte(year)[2]) + "0"

	geoParamValues := strings.Split(req.QueryStringParameters[geoParam], ",")

	fields := req.QueryStringParameters["fields"]

	for _, field := range strings.Split(fields, ",") {
		if matches, _ := regexp.MatchString(fieldPattern, field); !matches {
			return Response{StatusCode: 422}, errors.New("invalid field: " + field)
		}
	}

	tableName := "acs5." + geoUnit + "_state_" + subject + "_" + year

	fieldColumns := geoParam + ", " + fields

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ANY($1)", fieldColumns, tableName, geoParam)

	rows, err := db.Query(query, pq.Array(geoParamValues))
	if err != nil {
		return Response{StatusCode: 500}, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	var returnJSON []interface{}


	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range columns {
			values[i] = new(sql.RawBytes)
		}
		if err = rows.Scan(values...); err != nil {
			return Response{StatusCode: 500}, err
		}
		rowJSON := gabs.New()
		for i, value := range values {
			var val []byte
			val = *value.(*sql.RawBytes)
			if i == 0 {
				if _, err = rowJSON.Set(string(val), geoParam); err != nil {
					return Response{StatusCode: 500}, err
				}
			} else {
				intValue, err := strconv.Atoi(string(val))
				if err != nil {
					return Response{StatusCode: 500}, err
				}
				if _, err = rowJSON.Set(intValue, columns[i]); err != nil {
					return Response{StatusCode: 500}, err
				}
			}
		}
		returnJSON = append(returnJSON, rowJSON.Data())
	}

	jsonB, err := json.Marshal(returnJSON)
	if err != nil {
		return Response{StatusCode: 500}, err
	}

	// output is json data structure
	return Response{StatusCode: 200, Headers: headers, Body: string(jsonB)}, nil
}

func main() {
	lambda.Start(Handler)
}
