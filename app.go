package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"route-switcher-go/routeswitcher"
	"route-switcher-go/ruleservice"
)

func main() {
	test()
	s, err := ruleservice.NewRuleService("rules.json")
	if err != nil {
		panic(err)
	}
	rules := s.GetRules()

	fmt.Println("Rules:")
	for _, rule := range rules {
		fmt.Println(rule)
	}

	http.HandleFunc("/", routeswitcher.NewProxyHandler())
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func test() {
	// Example of creating a Rule and serializing it to JSON.
	rule := ruleservice.NewRule("example.com", "192.168.1.1", "opt1", "target1", "admin", 1623456000, nil)
	jsonBytes, _ := json.MarshalIndent(rule, "", "  ")
	fmt.Println(string(jsonBytes))

	// Demonstrating deserialization, note that handling of unknown fields being ignored is default behavior in Go's json package.
	var deserializedRule ruleservice.Rule
	if err := json.Unmarshal(jsonBytes, &deserializedRule); err != nil {
		fmt.Println("Error deserializing:", err)
	} else {
		fmt.Println("Deserialized Rule:", deserializedRule)
	}
}
