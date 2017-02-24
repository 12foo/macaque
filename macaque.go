package macaque

import (
	"github.com/mholt/caddy/caddyhttp/httpserver"
	"github.com/patrickmn/go-cache"
	"net/http"
	"time"
)

type MacaquePolicy struct {
	Name           string
	PerIPMax       int
	PerIPInterval  int
	PerKeyMax      int
	PerKeyInterval int
}

type MacaqueConfig struct {
	Path         string
	Policies     map[string]MacaquePolicy
	SyncInterval int
	DBPath       string
	PortalAPI    string
}

type Macaque struct {
	Next   httpserver.Handler
	Config MacaqueConfig
	mem    *cache.Cache
}

func (mq Macaque) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	return m.Next.ServeHTTP(w, r)
}

func NewMacaque(config MacaqueConfig) Macaque {
	mq := Macaque{Config: config}
	mq.mem = cache.New(mq.Config.SyncInterval*time.Minute, time.Minute)
	return mq
}
