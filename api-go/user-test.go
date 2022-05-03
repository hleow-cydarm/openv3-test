package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

func main() {
	ctx := context.Background()
	loader := openapi3.Loader{Context: ctx}
	doc, err := loader.LoadFromFile("../spec/user.yaml")
	if err != nil {
		panic(err)
	}
	err = doc.Validate(ctx)
	if err != nil {
		panic(err)
	}
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		panic(err)
	}
	getUsernamePath(router, ctx)
	getUserPath(router, ctx)
	postUserPath(router, ctx)
	invalidPostUserPath(router, ctx)

}

func getUsernamePath(router routers.Router, ctx context.Context) {
	httpReq, err := http.NewRequest(http.MethodGet, "/user/"+"username", nil)
	if err != nil {
		panic(err)
	}

	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}
	log.Println("Route:", route.Path)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	type User struct {
		Username string `json:"username" gorm:"username"`
		Name     string `json:"name" gorm:"name"`
	}
	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = &User{Name: "testUser"}
	)

	log.Println("Response:", respStatus)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	data, _ := json.Marshal(respBody)
	responseValidationInput.SetBodyBytes(data)
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}

	// changed response to invalid schema
	data, _ = json.Marshal([]*User{{Name: "testUser"}})
	responseValidationInput.SetBodyBytes(data)
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err == nil {
		log.Fatalf("expecting err over here as response body doesnt match the schema")
	}
	log.Println("Passed Username path")
}

func getUserPath(router routers.Router, ctx context.Context) {
	httpReq, err := http.NewRequest(http.MethodGet, "/user", nil)
	if err != nil {
		panic(err)
	}

	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}
	log.Println("Route:", route.Path)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		panic(err)
	}

	type User struct {
		Username string `json:"username" gorm:"username"`
		Name     string `json:"name" gorm:"name"`
	}
	var (
		respStatus      = 200
		respContentType = "application/json"
		respBody        = []*User{{Name: "a@a.com"}}
	)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	data, _ := json.Marshal(respBody)
	log.Println("Data: ", string(data))
	responseValidationInput.SetBodyBytes(data)
	// Validate response.
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
	// changed response to invalid schema
	data, _ = json.Marshal(User{Name: "testUser"})
	responseValidationInput.SetBodyBytes(data)
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err == nil {
		log.Fatalf("expecting err over here as response body doesnt match the schema")
	}
	log.Println("Passed User Path")
}

func postUserPath(router routers.Router, ctx context.Context) {

	type DraftUser struct {
		Name     string `json:"name" gorm:"name"`
		Password string `json:"password"`
	}

	reqBody := &DraftUser{Password: "a@a.com"}
	jsonBody, err := json.Marshal(&reqBody)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}
	log.Println("Route:", route.Path)

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
		respStatus      = 201
		respContentType = "application/json"
	)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{respContentType}},
	}
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		panic(err)
	}
	log.Println("Passed Post user path")
}

func invalidPostUserPath(router routers.Router, ctx context.Context) {
	//name is required in the path; in this test case we will not add name
	type DraftWrongSchemaUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	reqBody := &DraftWrongSchemaUser{Email: "a@a.com", Password: "pass"}
	jsonBody, err := json.Marshal(&reqBody)
	if err != nil {
		panic(err)
	}
	httpReq, err := http.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	// Find route
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		panic(err)
	}
	log.Println("Route:", route.Path)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err == nil {
		log.Fatalf("expecting err over here as request body doesnt match the schema; `Name` field is required")
	}
	log.Println("Passed invalidPostUser path")
}
