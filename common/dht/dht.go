package dht

import "github.com/sakshamsharma/sarga/common/network"

// DHT is a common interface to be satisfied by
// all implementations to be used with sarga.
type DHT interface {
	Init(seeds []network.Address) error
	FindValue(key string) ([]byte, error)
	StoreValue(key string, data []byte) error
	ShutDown() error
}
