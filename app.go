package main

import (
	"crypto/tls"
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

func startHttpsServer() {
	port := flag.Int("p", 0, "Port number to use, default is 0 for random")
	flag.Parse()

	cert, err := tls.LoadX509KeyPair("tls/fullchain.pem", "tls/privkey.pem")
	//cert, err := tls.LoadX509KeyPair("/root/tls/star-20230101-fullchain.pem", "/root/tls/star-20230101-privkey.pem")
	if err != nil {
		panic(err)
	}
	log.Println("TLS loaded")

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%d", *port), config)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	actPort := listener.Addr().(*net.TCPAddr).Port
	log.Printf("Server started at %d", actPort)

	err = http.Serve(listener, nil)
	if err != nil {
		panic(err)
	}
}
