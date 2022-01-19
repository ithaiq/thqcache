package thqcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_thqcache/"

type HTTPPoll struct {
	self     string
	basePath string
}

func NewHTTPPoll(self string) *HTTPPoll {
	return &HTTPPoll{
		self: self, basePath: defaultBasePath,
	}
}

func (this *HTTPPoll) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
}

func (this *HTTPPoll) Logf(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", this.self, fmt.Sprintf(format, v...))
}
