package handler

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"route-switcher-go/ruleservice"

	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type dynamicProxyHandler struct {
	ruleService           ruleservice.RuleService
	proxyHandler          http.Handler
	embeddedStaticHandler http.Handler
	rootStaticHandler     http.Handler
}

func NewProxyHandler(rs ruleservice.RuleService, dh http.Handler) http.Handler {
	switcher := dynamicProxyHandler{ruleService: rs, embeddedStaticHandler: dh}
	switcher.initProxy()
	return &switcher
}

func (it *dynamicProxyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	//log.Printf("inURL: %s, inUrlPath: %s, remoteAddr: %s", r.URL, r.URL.Path, r.RemoteAddr)
	rule := it.ruleService.FindRule(r.URL.Path, r.RemoteAddr)
	if rule.Target == "" {
		//log.Printf("No rule found for %s, %s", r.URL.Path, r.RemoteAddr)
		if r.URL.Path == "/route-switcher.log" || r.URL.Path == "/rules.json" {
			it.rootStaticHandler.ServeHTTP(rw, r)
		} else {
			it.embeddedStaticHandler.ServeHTTP(rw, r)
		}
	} else {
		uid, _ := uuid.NewV7()
		corrIdStr := uid.String()[len(uid.String())-12:]

		log.Printf("request:  [%s][%s][%s][%s] ===> %s%s", corrIdStr, r.URL.Path, r.RemoteAddr, r.Method, rule.Target, r.URL.Path)
		request := r.WithContext(context.WithValue(r.Context(), "correlationId", corrIdStr))
		it.proxyHandler.ServeHTTP(rw, request) // let the ReverseProxy to handle the request
	}
}

func (it *dynamicProxyHandler) initProxy() {
	it.proxyHandler = &httputil.ReverseProxy{Rewrite: it.rewrite, ModifyResponse: it.modifyResponse}
	it.rootStaticHandler = http.FileServer(http.Dir("./"))
}

func (it *dynamicProxyHandler) rewrite(r *httputil.ProxyRequest) {
	rule := it.ruleService.FindRule(r.In.URL.Path, r.In.RemoteAddr)
	target, _ := url.Parse(rule.Target)
	r.SetURL(target)
}

func (it *dynamicProxyHandler) modifyResponse(resp *http.Response) error {
	log.Printf("response: [%s][%s]", resp.Request.Context().Value("correlationId"), resp.Status)
	return nil
}
