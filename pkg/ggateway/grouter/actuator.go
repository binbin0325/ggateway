package grouter

import "C"
import (
	"bytes"
	"fmt"
	"ggateway/pkg/ggateway"
	"ggateway/pkg/rc/nacos"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/valyala/fasthttp"
	"io"
	"strconv"
	"sync"
	"time"
)

func actuator(c *ggateway.Context) {
	for _, v := range routerMapping {
		proxyReq(v, c)
	}
}

func getInstance(uri string) *model.Instance {
	instance, err := nacos.GetRegisterClient().SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: uri,
	})
	if err != nil {
		fmt.Println(err)
	}
	return instance

}

var bufferPool = bufferPoolFunc()

type Adapter struct {
	pool sync.Pool
}

func bufferPoolFunc() *Adapter {
	return &Adapter{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 4096))
			},
		},
	}
}

func (api *Adapter) getRequestBodyBytes(body io.ReadCloser) []byte {
	buffer := api.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			api.pool.Put(buffer)
			buffer = nil
		}
	}()

	_, err := io.Copy(buffer, body)
	if err != nil {
		fmt.Println(err)
	}
	return buffer.Bytes()
}

func proxyReq(v *Router, c *ggateway.Context) {
	defer fasthttp.ReleaseRequest(c.Req) // 用完需要释放资源
	var requestUrl string
	if v.Type == "lb" {
		instance := getInstance(v.Uri)
		if instance == nil {
			fmt.Println("instance is nil")
		}
		requestUrl = "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10) + c.Path
	} else {
		requestUrl = v.Uri
	}
	c.Req.SetRequestURI(requestUrl)
	c.Resp = fasthttp.AcquireResponse()
	err := fasthttp.DoTimeout(c.Req, c.Resp, 10*time.Second)

	if err != nil {
		fmt.Println(err)
	}
}
