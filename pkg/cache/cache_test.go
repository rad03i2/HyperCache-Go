package cache

import (
 "testing"
 "time"
)
func TestSetGetCopiesValue(t *testing.T){c:=New(Options{Shards:4,CleanupInterval:time.Hour});defer c.Close();v:=[]byte("hello");c.Set("a",v,0);v[0]='X';got,ok:=c.Get("a");if !ok||string(got)!="hello"{t.Fatalf("got %q",got)};got[0]='Y';again,_:=c.Get("a");if string(again)!="hello"{t.Fatal("Get exposed internal buffer")}}
func TestTTLExpiration(t *testing.T){c:=New(Options{Shards:2,CleanupInterval:5*time.Millisecond});defer c.Close();c.Set("x",[]byte("1"),15*time.Millisecond);time.Sleep(30*time.Millisecond);if _,ok:=c.Get("x");ok{t.Fatal("expired key returned")}}
func TestDeleteAndStats(t *testing.T){c:=New(Options{Shards:2,CleanupInterval:time.Hour});defer c.Close();c.Set("x",[]byte("1"),0);if !c.Delete("x"){t.Fatal("delete failed")};if c.Delete("x"){t.Fatal("second delete should fail")};_,_=c.Get("missing");s:=c.Snapshot();if s.Keys!=0||s.Sets!=1||s.Deletes!=1||s.Misses!=1{t.Fatalf("unexpected stats: %+v",s)}}
