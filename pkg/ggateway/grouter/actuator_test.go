package grouter

import (
	"fmt"
	"ggateway/pkg/ggateway"
	"ggateway/pkg/rc/nacos"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"testing"
)

func TestProxyReq(t *testing.T) {
	nacos.GetRegisterClient()
	r := &Router{Uri: "go_service_pre", Type: "lb"}
	c := &ggateway.Context{Req: fasthttp.AcquireRequest(), Resp: fasthttp.AcquireResponse(), Path: "/v1/log"}
	proxyReq(r, c)
}

func BenchmarkProxyReq(b *testing.B) {
	nacos.GetRegisterClient()
	r := &Router{Uri: "go_service_pre", Type: "lb"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := &ggateway.Context{Req: fasthttp.AcquireRequest(), Resp: fasthttp.AcquireResponse(), Path: "/v1/log"}
		proxyReq(r, c)

	}
}

func init() {
	viper.SetConfigName("config")                                //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("C:\\GoWork\\binbin\\ggateway\\configs") // 第一个搜索路径
	err := viper.ReadInConfig()                                  // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}
