package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"route-switcher-go/handler"
	"route-switcher-go/logging"
	"route-switcher-go/ruleservice"
)

func main() {
	logFile, logMsgChannel := logging.InitLogging()
	defer logFile.Close()

	ruleSvc, err := ruleservice.NewRuleService("rules.json")
	if err != nil {
		panic(err)
	}

	http.Handle("/log-msg", logging.NewLogMsgWebSocketHandler(logMsgChannel))
	http.Handle("/rule-manage/api/rules", ruleservice.NewRuleManageHandler(ruleSvc))
	http.Handle("/", handler.NewProxyHandler(ruleSvc, handler.NewEmbedStaticFileServer()))

	startServer()
}

func startServer() {
	port := flag.Int("p", 0, "Port number to use, default is 0 for random")
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprint(":", *port))
	if err != nil {
		panic(err)
	}
	actPort := listener.Addr().(*net.TCPAddr).Port
	log.Printf("Server started at %d", actPort)

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
