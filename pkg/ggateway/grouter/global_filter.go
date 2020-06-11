package grouter

import (
	"fmt"
	"ggateway/pkg/ggateway"
	"net/http"
)

func contextPathStripPrefixGlobalFilter() ggateway.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("123123",req.URL.Path)
	}
}

func contextPathStripPrefixGlobalFilter1() ggateway.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("456456",req.URL.Path)
	}
}

func loadGlobalFilters(router *ggateway.Router) {
	router.Use(ggateway.HandlerOrderFunc{
		Order:      2,
		FilterFunc: contextPathStripPrefixGlobalFilter1(),
	})

	router.Use(ggateway.HandlerOrderFunc{
		Order:      -1,
		FilterFunc: contextPathStripPrefixGlobalFilter(),
	})

}
