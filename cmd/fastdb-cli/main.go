// cmd/fastdb-server/main.go
package main

import (
	"bytes"
	"log"
	"runtime"

	"fastdb/internals/command"
	"fastdb/internals/engine"
	"fastdb/internals/protocol"

	"github.com/panjf2000/gnet"
)

type fastServer struct {
	*gnet.EventServer
	executor *command.Executor
}

type ConnCtx struct {
	id             string
	buf            *bytes.Buffer
	subscribedKeys map[string]struct{}
}

type ConnState struct {
	buffer         []byte
	id             string
	subscribedKeys map[string]struct{}
}

func newFastServer(ex *command.Executor) *fastServer {
	return &fastServer{executor: ex}
}

func (fs *fastServer) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
    state := &ConnState{
        buffer:         make([]byte, 0, 1024),
        id:             c.RemoteAddr().String(),
        subscribedKeys: make(map[string]struct{}),
    }
    c.SetContext(state)
    return
}

func (fs *fastServer) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
    ctxRaw := c.Context()
    if ctxRaw == nil {
        return
    }
    ctx := ctxRaw.(*ConnState)
    for key := range ctx.subscribedKeys {
        fs.executor.Engine.PubSub().Unsubscribe(key, ctx.id)
    }
    return
}

func (fs *fastServer) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	log.Printf("✅ FastDB server is running on %s (Multicore: %v)", srv.Addr.String(), srv.Multicore)
	return
}

func (fs *fastServer) React(packet []byte, c gnet.Conn) (out []byte, action gnet.Action) {
    state := c.Context().(*ConnState)

    state.buffer = append(state.buffer, packet...)

    for {
        args, consumed, err := protocol.ParseRESPCommandPartial(state.buffer)
        if err != nil {
            c.AsyncWrite(protocol.EncodeError(err.Error()))
            state.buffer = nil 
            return
        }

        if args == nil {
            break
        }

        state.buffer = state.buffer[consumed:]

        if len(args) == 0 {
            c.AsyncWrite(protocol.EncodeError("Empty command"))
            continue
        }

        cmd := args[0]

        response, err := fs.executor.ExecuteCommand(&command.Command{
            Name: cmd,
            Args: args[1:],
        })
        if err != nil {
            c.AsyncWrite(protocol.EncodeError(err.Error()))
            continue
        }

        if response == "" {
            c.AsyncWrite(protocol.EncodeEmpty)
        } else {
            c.AsyncWrite(protocol.EncodeString(response))
        }
    }

    return
}

func main() {
	db := engine.New(uint32(runtime.GOMAXPROCS(runtime.NumCPU())))
	ex := command.NewExecutor(db)
	server := newFastServer(ex)

	log.Println("🚀 Starting FastDB server on :6380")

	err := gnet.Serve(server, "tcp://:6380", gnet.WithMulticore(true))
	if err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
