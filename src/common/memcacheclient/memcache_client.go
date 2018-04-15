package memcacheclient


import (
	"errors"
	"fmt"
	"time"

	"github.com/gomemcache/memcache"
)

type MemcacheClient struct {
	Servers   []string
	TimeoutMs int

	Client *memcache.Client
}

func (client *MemcacheClient) Close() {
	if client.Client != nil {
		//todo, fix close
	}
}

func (client *MemcacheClient) Init() error {
	if len(client.Servers) == 0 || 0 > client.TimeoutMs {
		err := errors.New(fmt.Sprintf("invalid memcache config: servers %s timeout %d",
			client.Servers, client.TimeoutMs))
		return err
	}

	client.Client = memcache.New(client.Servers...)
	client.Client.Timeout = time.Duration(client.TimeoutMs) * time.Millisecond

	return nil
}

func (client *MemcacheClient) GetMulti(keys []string) (result map[string][]byte, err error) {
	if client.Client == nil {
		err = errors.New("internal error - MemcacheClient not init")
		return
	}

	data, err := client.Client.GetMulti(keys)
	if err != nil {
		return
	}

	result = make(map[string][]byte)
	for k, v := range data {
		result[k] = v.Value
	}
	return
}

func (client *MemcacheClient) Set(key string, data []byte) (err error) {
	if client.Client == nil {
		err = errors.New("internal error - MemcacheClient not init")
		return
	}

	item := &memcache.Item{Key: key, Value: data}
	err = client.Client.Set(item)
	return
}