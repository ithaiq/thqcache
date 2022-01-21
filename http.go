package thqcache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/ithaiq/thqcache/consistenthash"
	pb "github.com/ithaiq/thqcache/proto"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const (
	defaultBasePath = "/_thqcache/"
	defaultReplicas = 50
)

//HTTPPool HTTP服务端
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.RWMutex
	peers       *consistenthash.Map
	httpGetters map[string]*HttpGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self: self, basePath: defaultBasePath,
	}
}

func (this *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, this.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	this.Logf("%s %s", r.Method, r.URL.Path)
	parts := strings.SplitN(r.URL.Path[len(this.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	this.Logf("%s %s", groupName, key)

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no group cache: "+groupName, http.StatusNotFound)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := proto.Marshal(&pb.Response{Value: value.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

func (this *HTTPPool) Logf(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.self, fmt.Sprintf(format, v...))
}

func (this *HTTPPool) Set(peers ...string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.peers = consistenthash.New(defaultReplicas, nil)
	this.peers.Add(peers...)
	this.httpGetters = make(map[string]*HttpGetter, len(peers))
	for _, peer := range peers {
		this.httpGetters[peer] = &HttpGetter{baseUrl: peer + this.basePath}
	}
}

func (this *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	this.mu.Lock()
	defer this.mu.Unlock()
	peer := this.peers.Get(key)
	if peer != "" && peer != this.self {
		this.Logf("Pick peer %s", peer)
		return this.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

//HttpGetter 封装http获取缓存
type HttpGetter struct {
	baseUrl string
}

func (h HttpGetter) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v%v/%v",
		h.baseUrl,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	log.Println(u)
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}
	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}
	return nil
}

var _ PeerGetter = (*HttpGetter)(nil)
