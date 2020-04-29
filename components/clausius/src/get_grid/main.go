package main

import (
	"clausius/common"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	nb_rows, nb_cols, err := common.GetGridSize()
	if err != nil {
		panic(err)
	}

	store, err := common.ConnectToFunes()
	if err != nil {
		panic(err)
	}
	m, err := store.GetGrid(nb_rows, nb_cols)
	if err != nil {
		panic(err)
	}

	m["nb_rows"] = nb_rows
	m["nb_cols"] = nb_cols

	json, err := json.Marshal(m)

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{
		Body: string(json), StatusCode: 200, Headers: headers}, nil
}

func main() {
	lambda.Start(handle)
}
