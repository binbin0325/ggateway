package grouter

import (
	"encoding/json"
	"ggateway/pkg/cc"
	"ggateway/pkg/cc/nacos"
	"ggateway/pkg/ggateway"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Router struct {
	Uri        string      `json:"uri"`
	Type       string      `json:"type,omitempty"`
	Predicates []Predicate `json:"predicates,omitempty"`
	Order      int         `json:"order,omitempty"`
	Filters    []Filter    `json:"filters,omitempty"`
}

type Predicate struct {
	Args map[string]string `json:"args,omitempty"`
	Name string            `json:"name,omitempty"`
}

type Filter struct {
	Name string            `json:"name,omitempty"`
	Args map[string]string `json:"args,omitempty"`
}

//路由映射表
//key:ContextPath+pattern  不包括*
//value: *Router
var routerMapping map[string]*Router

var contextPath string

func InitRouter() (router *ggateway.Router) {
	contextPath = viper.GetString("server.router.context_path")
	router = ggateway.New()
	loadGlobalFilters(router)
	loadRouter(getRouters(), router)
	return
}

func getRouters() []*Router {
	var cc cc.ConfigCenter
	cc = &nacos.ConfigServer{
		Ip:   viper.GetString("nacos.config.ip"),
		Port: viper.GetUint64("nacos.config.port"),
	}
	configClient := cc.GetConfigClient().(config_client.ConfigClient)
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "routers",
		Group:  "gogateway",
	})
	if err != nil {
		os.Exit(1)
	}
	var routers []*Router
	err = json.Unmarshal([]byte(content), &routers)
	if err != nil {
		os.Exit(1)
	}
	return routers
}

func loadRouter(routers []*Router, router *ggateway.Router) {
	routerMapping = make(map[string]*Router)
	for _, r := range routers {
		for _, p := range r.Predicates {
			pattern := p.Args["pattern"]
			if index := strings.Index(pattern, "*"); index > 0 {
				requestKey := pattern[:index]
				routerMapping[contextPath+requestKey] = r
				pattern = pattern + "path"
			} else {
				routerMapping[contextPath+pattern] = r
			}
			router.Any(pattern, actuator)
		}
	}
}
