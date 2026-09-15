package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/rad03i2/HyperCache-Go/pkg/cache"
	"github.com/rad03i2/HyperCache-Go/pkg/resp"
)

type TCPServer struct {
	addr  string
	cache *cache.ShardedCache
}

func NewTCPServer(addr string, c *cache.ShardedCache) *TCPServer {
	return &TCPServer{
		addr:  addr,
		cache: c,
	}
}

func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("[HyperCache] Redis-compatible TCP server listening on %s\n", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := resp.NewReader(conn)
	writer := resp.NewWriter(conn)

	for {
		val, err := reader.ReadValue()
		if err != nil {
			if err != io.EOF {
				// Connection terminated
			}
			return
		}

		if val.Type != resp.TypeArray || len(val.Array) == 0 {
			writer.WriteError("expected command array")
			continue
		}

		cmd := strings.ToUpper(string(val.Array[0].Bulk))
		args := val.Array[1:]

		switch cmd {
		case "PING":
			writer.WriteSimpleString("PONG")

		case "SET":
			if len(args) < 2 {
				writer.WriteError("wrong number of arguments for 'set' command")
				continue
			}
			key := string(args[0].Bulk)
			value := args[1].Bulk
			s.cache.Set(key, value, 0)
			writer.WriteSimpleString("OK")

		case "GET":
			if len(args) < 1 {
				writer.WriteError("wrong number of arguments for 'get' command")
				continue
			}
			key := string(args[0].Bulk)
			val, found := s.cache.Get(key)
			if !found {
				writer.WriteBulk(nil)
			} else {
				writer.WriteBulk(val)
			}

		case "DEL":
			if len(args) < 1 {
				writer.WriteError("wrong number of arguments for 'del' command")
				continue
			}
			count := int64(0)
			for _, arg := range args {
				if s.cache.Delete(string(arg.Bulk)) {
					count++
				}
			}
			writer.WriteInteger(count)

		case "INFO":
			keys, hits, misses := s.cache.Stats()
			info := fmt.Sprintf("# HyperCache Stats\r\nkeys:%d\r\nhits:%d\r\nmisses:%d\r\nuptime:%d\r\n", keys, hits, misses, time.Now().Unix())
			writer.WriteBulk([]byte(info))

		default:
			writer.WriteError(fmt.Sprintf("unknown command '%s'", cmd))
		}
	}
}
