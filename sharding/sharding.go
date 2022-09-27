package sharding

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	m sync.RWMutex
	d map[string]interface{}
}

type ShardedMap []*Shard

func NewShardedMap(numberOfShards int) ShardedMap {
	shards := make([]*Shard, numberOfShards)

	for i := 0; i < numberOfShards; i++ {
		shards[i] = &Shard{d: make(map[string]interface{})}
	}

	return shards
}

func (sm ShardedMap) getShardIndex(key string) int {
	checksum := sha1.Sum([]byte(key))
	hash := int(checksum[17]) // any byte
	return hash % len(sm)
}

func (sm ShardedMap) getShard(key string) *Shard {
	index := sm.getShardIndex(key)
	return sm[index]
}

func (sm ShardedMap) Get(key string) interface{} {
	shard := sm.getShard(key)
	shard.m.RLock()
	defer shard.m.RUnlock()

	return shard.d[key]
}

func (sm ShardedMap) Set(key string, value interface{}) {
	shard := sm.getShard(key)
	shard.m.Lock()
	defer shard.m.Unlock()

	shard.d[key] = value
}
