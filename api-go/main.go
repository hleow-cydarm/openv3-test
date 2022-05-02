package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	legacyrouter "github.com/getkin/kin-openapi/routers/legacy"
)

func main() {

	myArray := [3][3]string{}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			myArray[i][j] = "."
		}
	}

	ctx := context.Background()
	loader := openapi3.Loader{Context: ctx}
	doc, _ := loader.LoadFromFile("../spec/tictactoe.yaml")
	_ = doc.Validate(ctx)
	router, _ := legacyrouter.NewRouter(doc)
	httpReq, _ := http.NewRequest(http.MethodGet, "/board", nil)

	// Find route
	route, pathParams, _ := router.FindRoute(httpReq)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = myArray
	)

	log.Println("Response:", respStatus)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	data, _ := json.Marshal(respBody)
	responseValidationInput.SetBodyBytes(data)

	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
}
