package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"route-switcher-go/handler"
	"route-switcher-go/ruleservice"
)

func main() {
	chanWriter := handler.NewChanWriter()
	logFile := initLogging(chanWriter)
	defer logFile.Close()

	ruleSvc, err := ruleservice.NewRuleService("rules.json")
	if err != nil {
		panic(err)
	}

	http.Handle("/log-msg", handler.NewLogMsgWebSocketHandler(chanWriter.Channel))
	http.Handle("/rule-manage/api/rules", ruleservice.NewRuleManageHandler(ruleSvc))
	http.Handle("/", handler.NewProxyHandler(ruleSvc, handler.NewEmbedStaticFileServer()))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func initLogging(chanWriter *handler.ChanWriter) *os.File {

	file, err := os.OpenFile("route-switcher.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	multiWriter := io.MultiWriter(os.Stdout, file, chanWriter)

	log.SetOutput(multiWriter)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("log initialized")
	return file
}
