package cache

import (
	"hash/fnv"
	"sync"
	"time"
)

type Item struct {
	Value      []byte
	ExpiresAt  int64 // Unix nanoseconds, 0 for no expiration
	LastAccess int64
}

func (i *Item) IsExpired() bool {
	if i.ExpiresAt == 0 {
		return false
	}
	return time.Now().UnixNano() > i.ExpiresAt
}

type CacheShard struct {
	mu    sync.RWMutex
	items map[string]*Item
}

type ShardedCache struct {
	numShards uint32
	shards    []*CacheShard
	hits      uint64
	misses    uint64
	muStats   sync.Mutex
}

func NewShardedCache(numShards uint32) *ShardedCache {
	if numShards == 0 {
		numShards = 64
	}
	sc := &ShardedCache{
		numShards: numShards,
		shards:    make([]*CacheShard, numShards),
	}
	for i := uint32(0); i < numShards; i++ {
		sc.shards[i] = &CacheShard{
			items: make(map[string]*Item),
		}
	}
	go sc.startEvictionLoop(time.Second * 5)
	return sc
}

func (sc *ShardedCache) getShard(key string) *CacheShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	shardIdx := h.Sum32() % sc.numShards
	return sc.shards[shardIdx]
}

func (sc *ShardedCache) Set(key string, value []byte, ttl time.Duration) {
	shard := sc.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	var exp int64 = 0
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	shard.items[key] = &Item{
		Value:      value,
		ExpiresAt:  exp,
		LastAccess: time.Now().UnixNano(),
	}
}

func (sc *ShardedCache) Get(key string) ([]byte, bool) {
	shard := sc.getShard(key)
	shard.mu.RLock()
	item, found := shard.items[key]
	shard.mu.RUnlock()

	if !found {
		sc.muStats.Lock()
		sc.misses++
		sc.muStats.Unlock()
		return nil, false
	}

	if item.IsExpired() {
		shard.mu.Lock()
		delete(shard.items, key)
		shard.mu.Unlock()
		sc.muStats.Lock()
		sc.misses++
		sc.muStats.Unlock()
		return nil, false
	}

	sc.muStats.Lock()
	sc.hits++
	sc.muStats.Unlock()
	return item.Value, true
}

func (sc *ShardedCache) Delete(key string) bool {
	shard := sc.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.items[key]; exists {
		delete(shard.items, key)
		return true
	}
	return false
}

func (sc *ShardedCache) startEvictionLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		now := time.Now().UnixNano()
		for _, shard := range sc.shards {
			shard.mu.Lock()
			for k, v := range shard.items {
				if v.ExpiresAt > 0 && now > v.ExpiresAt {
					delete(shard.items, k)
				}
			}
			shard.mu.Unlock()
		}
	}
}

func (sc *ShardedCache) Stats() (totalKeys int, hits uint64, misses uint64) {
	sc.muStats.Lock()
	hits = sc.hits
	misses = sc.misses
	sc.muStats.Unlock()

	for _, shard := range sc.shards {
		shard.mu.RLock()
		totalKeys += len(shard.items)
		shard.mu.RUnlock()
	}
	return
}
