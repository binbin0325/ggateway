package grouter

import (
	"encoding/json"
	"fmt"
	"ggateway/pkg/cc"
	"ggateway/pkg/cc/nacos"
	"ggateway/pkg/ggateway"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
)

type Router struct {
	Uri                string      `json:"uri"`
	Type               string      `json:"type,omitempty"`
	Predicates         []Predicate `json:"predicates,omitempty"`
	Order              int         `json:"order,omitempty"`
	Filters            []Filter    `json:"filters,omitempty"`
	RouterFiltersChain RouterFiltersChain
}

type Predicate struct {
	Args map[string]string `json:"args,omitempty"`
	Name string            `json:"name,omitempty"`
}

type Filter struct {
	Name string            `json:"name,omitempty"`
	Args map[string]string `json:"args,omitempty"`
}

// RouterFiltersChain defines a RouterFilterFunc array.
type RouterFiltersChain []RouterFilterOrder

//Router Filter Order Func
type RouterFilterOrder struct {
	order            int64
	routerFilterFunc RouterFilterFunc
}

// RouterFilterFunc defines the handler used by gin middleware as return value.
type RouterFilterFunc func(req *http.Request) error

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

func Index(w http.ResponseWriter, r *http.Request, _ ggateway.Params) {
	fmt.Println("hahahha")
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

func actuator(w http.ResponseWriter, req *http.Request,_ ggateway.Params) {
	for k, v := range routerMapping {
		fmt.Println(k, v)
	}
}
