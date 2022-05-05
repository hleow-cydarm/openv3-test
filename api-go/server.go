package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

type User struct {
	Username string `json:"username" gorm:"username"`
	Name     string `json:"name" gorm:"name"`
}

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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := validateFunc(router, ctx, r, w)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	})
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(route *routers.Route, httpReq *http.Request, pathParams map[string]string) (int, interface{}) {
	path := route.Path
	method := route.Method

	if path == "/user" && method == http.MethodPost {
		return http.StatusCreated, &User{Name: "created test user", Username: "created test userName"}
	}

	if path == "/user/{username}" && pathParams["username"] != "" && httpReq.Method == "GET" {
		return http.StatusOK, &User{Name: "found user", Username: "found username " + pathParams["username"]}
	}

	if path == "/user" && httpReq.Method == "GET" {
		return http.StatusOK, []*User{{Name: "user 1", Username: "user 1"}, {Name: "user 2", Username: "user 2"}}
	}
	return 500, nil
}

func validateFunc(router routers.Router, ctx context.Context, httpReq *http.Request, w http.ResponseWriter) error {
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("find route err: %v", err)
	}
	log.Println("Route:", route.Path)

	// Validate request
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	if err := openapi3filter.ValidateRequest(ctx, requestValidationInput); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return fmt.Errorf("validate request err: %v", err)
	}

	respStatus, respData := handleRequest(route, httpReq, pathParams)

	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 respStatus,
		Header:                 http.Header{"Content-Type": []string{"application/json"}},
	}
	data, err := json.Marshal(respData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("error parsing json %v", err)
	}
	responseValidationInput.SetBodyBytes(data)
	if err := openapi3filter.ValidateResponse(ctx, responseValidationInput); err != nil {
		return fmt.Errorf("error validate response: %v", err)
	}
	log.Println("Path checked")
	JsonResponse(w, respStatus, respData)
	return nil
}

func JsonResponse(w http.ResponseWriter, status int, jsonData interface{}) {
	if err := json.NewEncoder(w).Encode(jsonData); err != nil {
		fmt.Errorf("error parsing json %v", err)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
}
