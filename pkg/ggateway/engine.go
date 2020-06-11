package ggateway

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

type httpServer struct {
	*gnet.EventServer
	pool   *goroutine.Pool
	router *Router
}

type httpCodec struct {
}

var errMsg = "Internal Server Error"
var errMsgBytes = []byte(errMsg)

func (h httpCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	return buf, nil
}

func (h httpCodec) Decode(c gnet.Conn) (out []byte, err error) {
	buf := c.Read()
	c.ResetBuffer()
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	c.SetContext(Context{req: req, w: new(GatewayHTTPResponseWriter), index: -1})
	if len(buf) > 0 {
		return buf, err
	} else {
		return
	}
}

func Server(port int, multicore bool, router *Router) {
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
	return
}

func (hs *httpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	cpuProfile, _ := os.Create("cpu_profile")
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()
	if c.Context() == nil {
		// bad thing happened
		out = errMsgBytes
		action = gnet.Close
		return
	}
	context := c.Context().(Context)
	if context.req != nil {
		context.router = hs.router
		context.ServeHTTP(context.w, context.req)
	}
	out = frame
	return
}

type Context struct {
	req    *http.Request
	w      *GatewayHTTPResponseWriter
	index  int8
	router *Router
}
