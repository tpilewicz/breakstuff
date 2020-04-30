package main

import (
	"clausius/common"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strconv"
)

func main() {
	lambda.Start(handle)
}

func handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	nbRows, nbCols, err := common.GetGridSize()
	if err != nil {
		panic(err)
	}

	xStr, yStr, vStr, ok := get_params(request.QueryStringParameters)
	if !ok {
		panic(fmt.Errorf("It seems a query parameter is missing"))
	}
	x, y, v, err := convertParams(xStr, yStr, vStr, nbRows, nbCols)
	if err != nil {
		panic(err)
	}

	store, err := common.ConnectToFunes()
	if err != nil {
		panic(err)
	}
	err = store.SetCell(x, y, v)
	if err != nil {
		panic(err)
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"

	return events.APIGatewayProxyResponse{
		Body: "Done", StatusCode: 200, Headers: headers}, nil
}

func get_params(params map[string]string) (string, string, string, bool) {
	x, okx := params["x"]
	y, oky := params["y"]
	v, okv := params["v"]
	return x, y, v, (okx && oky && okv)
}

func convertParams(xStr string, yStr string, vStr string, nbRows int, nbCols int) (int, int, int, error) {
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return 0, 0, 0, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return 0, 0, 0, err
	}
	v, err := strconv.Atoi(vStr)
	if err != nil {
		return 0, 0, 0, err
	}

	if x < 0 || y < 0 || x >= nbCols || y >= nbRows {
		return x, y, v, fmt.Errorf("weird values for x and/or y. x: %v, y: %v", x, y)
	}
	if !(v == 0 || v == 1) {
		return x, y, v, fmt.Errorf("v must be 0 or 1. got: %v", v)
	}

	return x, y, v, nil
}
