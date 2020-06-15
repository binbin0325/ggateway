package grouter

import (
	"ggateway/pkg/ggateway"
	"regexp"
	"strings"
)
var regexpVersion *regexp.Regexp

func contextPathStripPrefixGlobalFilter() ggateway.HandlerFunc {
	return func(c *ggateway.Context) {
		c.Req.URI().SetPath(c.Path[len(contextPath):])
	}
}

func versionStripPrefixGlobalFilter() ggateway.HandlerFunc {
	return func(c *ggateway.Context) {
		versionRegxp := regexpVersion.FindStringSubmatch(c.Path)
		if len(versionRegxp) > 0 {
			if index := strings.Index(c.Path, versionRegxp[0]); index > 0 {
				c.Path = c.Path[index:]
			}
		}
	}
}

func loadGlobalFilters(router *ggateway.Router) {
	regexpVersion = regexp.MustCompile(`/v\d/`)
	router.Use(ggateway.HandlerOrderFunc{
		Order:      -100,
		FilterFunc: contextPathStripPrefixGlobalFilter(),
	})
	router.Use(ggateway.HandlerOrderFunc{
		Order:      -99,
		FilterFunc: versionStripPrefixGlobalFilter(),
	})
}
