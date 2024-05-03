package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"route-switcher-go/ruleservice"
)

type dynamicProxyHandler struct {
	ruleService    ruleservice.RuleService
	proxyHandler   http.Handler
	defaultHandler http.Handler
}

func NewProxyHandler(rs ruleservice.RuleService, dh http.Handler) http.Handler {
	switcher := dynamicProxyHandler{ruleService: rs, defaultHandler: dh}
	switcher.initProxy()
	return &switcher
}

func (it *dynamicProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("inURL: %s, inUrlPath: %s, remoteAddr: %s", r.URL, r.URL.Path, r.RemoteAddr)
	rule := it.ruleService.FindRule(r.URL.Path, r.RemoteAddr)
	if rule.Target == "" {
		log.Printf("no rule found for %s, %s", r.URL.Path, r.RemoteAddr)
		it.defaultHandler.ServeHTTP(w, r)
	} else {
		log.Printf("rule found: %s", rule)
		it.proxyHandler.ServeHTTP(w, r)
	}
}

func (it *dynamicProxyHandler) initProxy() {
	it.proxyHandler = &httputil.ReverseProxy{Rewrite: it.rewrite}
}

func (it *dynamicProxyHandler) rewrite(r *httputil.ProxyRequest) {
	rule := it.ruleService.FindRule(r.In.URL.Path, r.In.RemoteAddr)
	target, _ := url.Parse(rule.Target)
	r.SetURL(target)
}
