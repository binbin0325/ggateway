//Alibaba Nacos config center
package nacos

import (
	"fmt"
	"ggateway/pkg/cc"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/nacos_client"
	"github.com/nacos-group/nacos-sdk-go/common/http_agent"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var clientConfigTest = constant.ClientConfig{
	TimeoutMs:           10 * 1000,
	BeatInterval:        5 * 1000,
	ListenInterval:      300 * 1000,
	NotLoadCacheAtStart: true,
	//Username:            "nacos",
	//Password:            "nacos",
}
var once sync.Once
var configClient config_client.ConfigClient

type ConfigServer struct {
	Ip string
	Port uint64
	Username string
	Password string
}

func (cs *ConfigServer) GetConfigClient() config_client.ConfigClient {
	//实现单例
	once.Do(func() {
		configClient = cs.initConfigClientTest()
	})
	return configClient

}

func (cs *ConfigServer) initConfigClientTest() config_client.ConfigClient {
	nc := nacos_client.NacosClient{}
	nc.SetServerConfig([]constant.ServerConfig{constant.ServerConfig{
		IpAddr:      viper.GetString("nacos.config.ip"),
		Port:        viper.GetUint64("nacos.config.port"),
		ContextPath: "/nacos",
	}})
	nc.SetClientConfig(clientConfigTest)
	nc.SetHttpAgent(&http_agent.HttpAgent{})
	client, _ := config_client.NewConfigClient(&nc)
	return client
}

func main() {
	client := initConfigClientTest()
	content, _ := client.GetConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	fmt.Println("config :" + content)
	_, err := client.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: "hello world!"})
	if err != nil {
		fmt.Printf("success err:%s", err.Error())
	}
	content = ""

	client.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data)
			content = data
		},
	})

	client.ListenConfig(vo.ConfigParam{
		DataId: "abc",
		Group:  "DEFAULT_GROUP",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

	time.Sleep(5 * time.Second)
	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: "abc"})

	select {}

}
