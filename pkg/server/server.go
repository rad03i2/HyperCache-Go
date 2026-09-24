package server

import (
 "fmt"
 "io"
 "log"
 "net"
 "strconv"
 "strings"
 "time"
 "github.com/rad03i2/HyperCache-Go/pkg/cache"
 "github.com/rad03i2/HyperCache-Go/pkg/resp"
)

type TCPServer struct { addr string; cache *cache.ShardedCache; started time.Time }
func NewTCPServer(addr string,c *cache.ShardedCache)*TCPServer{return &TCPServer{addr:addr,cache:c,started:time.Now()}}
func (s *TCPServer) Start() error { l,e:=net.Listen("tcp",s.addr);if e!=nil{return e};defer l.Close();log.Printf("HyperCache listening on %s",s.addr);for{c,e:=l.Accept();if e!=nil{return e};go s.handleConnection(c)} }
func arg(v resp.Value) string { if v.Type==resp.TypeBulkString{return string(v.Bulk)};return v.Str }
func (s *TCPServer) handleConnection(c net.Conn){defer c.Close();r:=resp.NewReader(c);w:=resp.NewWriter(c);for{v,e:=r.ReadValue();if e!=nil{if e!=io.EOF{_ = w.WriteError("protocol error")};return};if v.Type!=resp.TypeArray||len(v.Array)==0{_ = w.WriteError("expected command array");continue};cmd:=strings.ToUpper(arg(v.Array[0]));a:=v.Array[1:];switch cmd{
case "PING": if len(a)==0{_ = w.WriteSimpleString("PONG")}else if len(a)==1{_ = w.WriteBulk([]byte(arg(a[0])))}else{_ = w.WriteError("wrong number of arguments for 'ping'")}
case "SET": if len(a)!=2&&len(a)!=4{_ = w.WriteError("wrong number of arguments for 'set'");continue};ttl:=time.Duration(0);if len(a)==4{n,e:=strconv.ParseInt(arg(a[3]),10,64);opt:=strings.ToUpper(arg(a[2]));if e!=nil||n<=0||(opt!="EX"&&opt!="PX"){_ = w.WriteError("invalid expire option");continue};if opt=="EX"{ttl=time.Duration(n)*time.Second}else{ttl=time.Duration(n)*time.Millisecond}};s.cache.Set(arg(a[0]),[]byte(arg(a[1])),ttl);_ = w.WriteSimpleString("OK")
case "GET": if len(a)!=1{_ = w.WriteError("wrong number of arguments for 'get'");continue};b,ok:=s.cache.Get(arg(a[0]));if !ok{_ = w.WriteBulk(nil)}else{_ = w.WriteBulk(b)}
case "DEL": if len(a)<1{_ = w.WriteError("wrong number of arguments for 'del'");continue};var n int64;for _,x:=range a{if s.cache.Delete(arg(x)){n++}};_ = w.WriteInteger(n)
case "EXISTS": if len(a)!=1{_ = w.WriteError("wrong number of arguments for 'exists'");continue};if s.cache.Exists(arg(a[0])){_ = w.WriteInteger(1)}else{_ = w.WriteInteger(0)}
case "TTL": if len(a)!=1{_ = w.WriteError("wrong number of arguments for 'ttl'");continue};d,ok:=s.cache.TTL(arg(a[0]));if !ok{_ = w.WriteInteger(-2)}else if d<0{_ = w.WriteInteger(-1)}else{_ = w.WriteInteger(int64(d/time.Second))}
case "DBSIZE": st:=s.cache.Snapshot();_ = w.WriteInteger(int64(st.Keys))
case "INFO": st:=s.cache.Snapshot();info:=fmt.Sprintf("# HyperCache\r\nkeys:%d\r\nhits:%d\r\nmisses:%d\r\nsets:%d\r\ndeletes:%d\r\nuptime_seconds:%d\r\n",st.Keys,st.Hits,st.Misses,st.Sets,st.Deletes,int64(time.Since(s.started).Seconds()));_ = w.WriteBulk([]byte(info))
case "QUIT": _ = w.WriteSimpleString("OK");return
default:_ = w.WriteError(fmt.Sprintf("unknown command '%s'",cmd))}}}
