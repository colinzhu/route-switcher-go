package main

import (
	"net/http"
	"route-switcher-go/ruleservice"
)

func main() {
	ruleSvc, err := ruleservice.NewRuleService("rules.json")
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/rule-manage/api/rules", ruleservice.NewRuleManageHandler(ruleSvc))
	http.Handle("/", NewProxyHandler(ruleSvc, fs))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
