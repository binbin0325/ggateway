package filter

import (
	"ggateway/pkg/ggateway"
)

func contextPathStripPrefixGlobalFilter() ggateway.HandlerFunc {
	return func(c *ggateway.Context) {
	}
}

func contextPathStripPrefixGlobalFilter1() ggateway.HandlerFunc {
	return func(c *ggateway.Context) {
	}
}

func LoadGlobalFilters(router *ggateway.Router) {
	router.Use(ggateway.HandlerOrderFunc{
		Order:      2,
		FilterFunc: contextPathStripPrefixGlobalFilter1(),
	})

	router.Use(ggateway.HandlerOrderFunc{
		Order:      -1,
		FilterFunc: contextPathStripPrefixGlobalFilter(),
	})

}
