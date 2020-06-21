package ggateway

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/spf13/viper"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

type Context struct {
	index  int8
	router *Router
	Resp   *fasthttp.Response
	Req    *fasthttp.Request
	Ps     Params
	Code   int
	Path   string
}
type httpServer struct {
	*gnet.EventServer
	pool   *goroutine.Pool
	router *Router
}

type httpCodec struct {
}

var errMsg = "Internal Server Error"
var errMsgBytes = []byte(errMsg)

var cpuProfile *os.File

func (h httpCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (h httpCodec) Decode(c gnet.Conn) (out []byte, err error) {
	buf := c.Read()
	c.ResetBuffer()
	if len(buf) > 0 {
		req := new(fasthttp.Request)
		req.Read(bufio.NewReader(bytes.NewReader(buf)))
		c.SetContext(Context{Req: req, index: -1})
		return buf, err
	} else {
		return
	}
}

func Server(port int, multicore bool, router *Router) {
	cpuProfile, _ = os.Create("cpu_profile")

	p := goroutine.Default()
	defer p.Release()

	http := &httpServer{pool: p, router: router}
	hc := new(httpCodec)
	// Start serving!
	log.Fatal(gnet.Serve(http, fmt.Sprintf("tcp://:%d", port), gnet.WithMulticore(multicore), gnet.WithCodec(hc)))
}

func (hs *httpServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("HTTP server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	hs.router.SortGlobalFilters()
	go func() {
		log.Println(http.ListenAndServe(viper.GetString("pprof.server"), nil))
	}()
	return
}

func (hs *httpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if c.Context() == nil {
		// bad thing happened
		out = errMsgBytes
		action = gnet.Close
		return
	}
	// Use ants pool to unblock the event-loop.
	_ = hs.pool.Submit(func() {
		ctx := c.Context().(Context)
		if ctx.Req != nil {
			ctx.router = hs.router
			ctx.ServeHTTP()
			if ctx.Resp != nil {
				if out, err := writerResp(&ctx); err != nil {
					fmt.Println(err)
				} else {
					c.AsyncWrite(out)
				}
			}
		}
	})

	return
}

func writerResp(ctx *Context) (out []byte, err error) {
	defer fasthttp.ReleaseResponse(ctx.Resp) // 用完需要释放资源
	buffer := bytebufferpool.Get()
	bw := bufio.NewWriter(buffer)
	err = ctx.Resp.Write(bw)
	err = bw.Flush()
	out = buffer.Bytes()
	return
}
