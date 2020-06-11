package grouter

import (
	"bufio"
	"bytes"
	"fmt"
	"ggateway/pkg/ggateway"
	"ggateway/pkg/rc/nacos"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func actuator(w http.ResponseWriter, req *http.Request, _ ggateway.Params) {
	for k, v := range routerMapping {
		fmt.Println(k, v)
		proxyReq(v,req,w)
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

func proxyReq(v *Router, req *http.Request,w http.ResponseWriter) {
	var requestUrl string
	if v.Type == "lb" {
		instance := getInstance(v.Uri)
		requestUrl = "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10) + req.URL.Path
	} else {
		requestUrl = v.Uri
	}

	fastReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(fastReq) // 用完需要释放资源
	// 默认是application/x-www-form-urlencoded
	fastReq.Header.SetContentType("application/json")
	fastReq.Header.SetMethod(req.Method)
	fastReq.SetRequestURI(requestUrl)
	fastReq.SetBody(bufferPool.getRequestBodyBytes(req.Body))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp) // 用完需要释放资源

	err := fasthttp.Do(fastReq, resp)
	if err!=nil{
		fmt.Println(err)
	}
	buffer :=bufferPool.pool.Get().(*bytes.Buffer)
	bw := bufio.NewWriter(buffer)
	if err := resp.Write(bw); err != nil {
		fmt.Println(err)
	}
	bw.Flush()
	w.Write(buffer.Bytes())
}
