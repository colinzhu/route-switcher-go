package routeswitcher

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type routeSwitcher struct {
	rules map[string]string
}

func NewProxyHandler() func(http.ResponseWriter, *http.Request) {
	switcher := routeSwitcher{}
	return switcher.proxyHandler()
}

func (it *routeSwitcher) proxyHandler() func(http.ResponseWriter, *http.Request) {
	proxy := &httputil.ReverseProxy{Rewrite: it.rewrite}
	return proxy.ServeHTTP
}

func (it *routeSwitcher) rewrite(r *httputil.ProxyRequest) {
	target, _ := it.getTarget(r)
	r.SetURL(target)
}

func (it *routeSwitcher) getTarget(r *httputil.ProxyRequest) (*url.URL, error) {
	log.Printf("get target %s", it.rules)
	return url.Parse("https://httpbin.org")
}
