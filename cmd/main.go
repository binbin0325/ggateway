package main

import (
	"fmt"
	"ggateway/pkg/ggateway"
	"ggateway/pkg/ggateway/grouter"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
)

var (
	commit  = "N/A"
	date    = "N/A"
	version = "N/A"
)

func main() {
	var (
		port      = kingpin.Flag("port", "Server Port.").Default(viper.GetString("server.port")).Int()
		multicore = kingpin.Flag("multicore", "Multi Core").Bool()
	)
	kingpin.Parse()
	quit := make(chan struct{})
	defer close(quit)
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)
	go func() {
		for {
			select { //nolint: megacheck
			case <-sig:
				quit <- struct{}{}
			}
			os.Exit(0)
		}
	}()
	ggateway.Server(*port, *multicore,grouter.InitRouter())
}

func init(){
	viper.SetConfigName("config")  //  设置配置文件名 (不带后缀)
	viper.AddConfigPath("configs") // 第一个搜索路径
	err := viper.ReadInConfig()    // 搜索路径，并读取配置数据
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:",  e.Name)
	})
}