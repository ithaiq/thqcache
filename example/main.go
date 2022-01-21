package main

import (
	"flag"
	"fmt"
	"github.com/ithaiq/thqcache"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var db = map[string]string{
	"thq1": "1",
	"thq2": "2",
	"thq3": "3",
}

func createGroup() *thqcache.Group {
	return thqcache.NewGroup("name", 2<<10, thqcache.GetterFunc(
		func(key string) ([]byte, error) {
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(addr string, addrs []string, group *thqcache.Group) {
	peers := thqcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Println("thqcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(apiAddr string, group *thqcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := group.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	var thqf string
	var thqs string
	flag.IntVar(&port, "port", 8001, "thqcache server port")
	flag.BoolVar(&api, "api", false, "start a api server")
	flag.StringVar(&thqf, "thq_f", "thq4", "")
	flag.StringVar(&thqs, "thq_s", "thq5", "")
	flag.Parse()

	rand.Seed(time.Now().Unix())
	db[thqf] = strconv.Itoa(rand.Intn(100))
	db[thqs] = strconv.Itoa(rand.Intn(100))
	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	group := createGroup()
	if api {
		go startAPIServer(apiAddr, group)
	}
	startCacheServer(addrMap[port], addrs, group)
}
