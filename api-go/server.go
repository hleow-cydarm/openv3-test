package main

import (
	"context"
	"fmt"
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
	log.Println("Path checked")
	return nil
}
