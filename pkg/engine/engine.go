package engine

import (
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"github.com/valyala/fasthttp"
	"log"
)

type httpServer struct {
	*gnet.EventServer
	pool *goroutine.Pool
}

type httpCodec struct {
	req fasthttp.Request
}

var errMsg = "Internal Server Error"
var errMsgBytes = []byte(errMsg)

func (h httpCodec) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	panic("implement me")
}

func (h httpCodec) Decode(c gnet.Conn) ([]byte, error) {
	panic("implement me")
}

func Server(port int, multicore bool) {
	p := goroutine.Default()
	defer p.Release()

	http := &httpServer{pool: p}
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
	req := c.Context().(fasthttp.Request)
	fmt.Println(string(req.RequestURI()))
	/*	data := append([]byte{}, frame...)

		// Use ants pool to unblock the event-loop.
		_ = hs.pool.Submit(func() {
			time.Sleep(1 * time.Second)
			c.AsyncWrite(data)
		})*/
	out = frame
	return
}