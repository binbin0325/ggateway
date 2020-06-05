package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"ggateway/pkg/ggateway"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"log"
	"net/http"
	_ "net/http/pprof"
)

type httpServer struct {
	*gnet.EventServer
	pool   *goroutine.Pool
	router *ggateway.Router
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
	c.SetContext(Context{req: req, w: new(ggateway.GatewayHTTPResponseWriter)})
	if len(buf) > 0 {
		return buf, err
	} else {
		return
	}
}

func Server(port int, multicore bool, router *ggateway.Router) {
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
	return
}

func (hs *httpServer) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if c.Context() == nil {
		// bad thing happened
		out = errMsgBytes
		action = gnet.Close
		return
	}
	context := c.Context().(Context)
	if context.req != nil {
		hs.router.ServeHTTP(context.w, context.req)
	}

	/*	data := append([]byte{}, frame...)

		// Use ants pool to unblock the event-loop.
		_ = hs.pool.Submit(func() {
			time.Sleep(1 * time.Second)
			c.AsyncWrite(data)
		})*/
	out = frame
	return
}



type Context struct {
	req *http.Request
	w   *ggateway.GatewayHTTPResponseWriter
}
