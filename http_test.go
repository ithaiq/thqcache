package thqcache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func HTTPServer(t *testing.T) {
	NewGroup("names", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	addr := "localhost:8888"
	peers := NewHTTPPoll(addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
