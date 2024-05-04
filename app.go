package main

import (
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

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
