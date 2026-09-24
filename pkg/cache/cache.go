package cache

import (
	"hash/fnv"
	"sync"
	"sync/atomic"
	"time"
)

type Item struct { Value []byte; ExpiresAt int64 }
func (i Item) expired(now int64) bool { return i.ExpiresAt > 0 && now >= i.ExpiresAt }

type shard struct { mu sync.RWMutex; items map[string]Item }
type Options struct { Shards uint32; CleanupInterval time.Duration }
type Stats struct { Keys int `json:"keys"`; Hits uint64 `json:"hits"`; Misses uint64 `json:"misses"`; Sets uint64 `json:"sets"`; Deletes uint64 `json:"deletes"` }

type ShardedCache struct {
	shards []*shard
	hits atomic.Uint64; misses atomic.Uint64; sets atomic.Uint64; deletes atomic.Uint64
	stop chan struct{}; done chan struct{}; once sync.Once
}

func NewShardedCache(n uint32) *ShardedCache { return New(Options{Shards:n, CleanupInterval:5*time.Second}) }
func New(o Options) *ShardedCache {
	if o.Shards == 0 { o.Shards = 64 }; if o.CleanupInterval <= 0 { o.CleanupInterval = 5*time.Second }
	c := &ShardedCache{shards:make([]*shard,o.Shards), stop:make(chan struct{}), done:make(chan struct{})}
	for i := range c.shards { c.shards[i] = &shard{items:make(map[string]Item)} }
	go c.cleanup(o.CleanupInterval); return c
}
func (c *ShardedCache) shardFor(key string) *shard { h:=fnv.New32a(); _,_=h.Write([]byte(key)); return c.shards[h.Sum32()%uint32(len(c.shards))] }
func (c *ShardedCache) Set(key string, value []byte, ttl time.Duration) {
	if key=="" { return }; var exp int64; if ttl>0 { exp=time.Now().Add(ttl).UnixNano() }
	v:=append([]byte(nil),value...); s:=c.shardFor(key); s.mu.Lock(); s.items[key]=Item{Value:v,ExpiresAt:exp}; s.mu.Unlock(); c.sets.Add(1)
}
func (c *ShardedCache) Get(key string) ([]byte,bool) {
	s:=c.shardFor(key); s.mu.RLock(); item,ok:=s.items[key]; s.mu.RUnlock()
	if !ok { c.misses.Add(1); return nil,false }
	if item.expired(time.Now().UnixNano()) { s.mu.Lock(); if cur,exists:=s.items[key]; exists && cur.expired(time.Now().UnixNano()) { delete(s.items,key) }; s.mu.Unlock(); c.misses.Add(1); return nil,false }
	c.hits.Add(1); return append([]byte(nil),item.Value...),true
}
func (c *ShardedCache) Delete(key string) bool { s:=c.shardFor(key); s.mu.Lock(); _,ok:=s.items[key]; if ok { delete(s.items,key) }; s.mu.Unlock(); if ok { c.deletes.Add(1) }; return ok }
func (c *ShardedCache) Exists(key string) bool { _,ok:=c.Get(key); return ok }
func (c *ShardedCache) TTL(key string) (time.Duration,bool) { s:=c.shardFor(key); s.mu.RLock(); i,ok:=s.items[key]; s.mu.RUnlock(); if !ok||i.expired(time.Now().UnixNano()){return 0,false}; if i.ExpiresAt==0{return -1,true}; return time.Until(time.Unix(0,i.ExpiresAt)),true }
func (c *ShardedCache) Snapshot() Stats { st:=Stats{Hits:c.hits.Load(),Misses:c.misses.Load(),Sets:c.sets.Load(),Deletes:c.deletes.Load()}; now:=time.Now().UnixNano(); for _,s:=range c.shards{s.mu.RLock();for _,i:=range s.items{if !i.expired(now){st.Keys++}};s.mu.RUnlock()};return st }
func (c *ShardedCache) Stats()(int,uint64,uint64){s:=c.Snapshot();return s.Keys,s.Hits,s.Misses}
func (c *ShardedCache) cleanup(d time.Duration){defer close(c.done); t:=time.NewTicker(d);defer t.Stop();for{select{case now:=<-t.C:n:=now.UnixNano();for _,s:=range c.shards{s.mu.Lock();for k,i:=range s.items{if i.expired(n){delete(s.items,k)}};s.mu.Unlock()};case <-c.stop:return}}}
func (c *ShardedCache) Close(){c.once.Do(func(){close(c.stop);<-c.done})}
