package grouter

import (
	"fmt"
	"ggateway/pkg/ggateway"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"
	"net/http"
	"os"
)

func InitRouter() *ggateway.Router {
	router := ggateway.New()
	router.Any("/", Index)
	return router
}

func Index(w http.ResponseWriter, r *http.Request, _ ggateway.Params) {
	fmt.Println("hahahha")
}

func getRouterContent() string {
	client := config.GetConfigClient()
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "routers",
		Group:  "gogateway",
	})
	if err != nil {
		log.Log.Error("err", zap.Error(err))
		os.Exit(1)
	}
	return content
}
