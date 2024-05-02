package main

import (
	"net/http"
	"route-switcher-go/routeswitcher"
)

func main() {
	http.HandleFunc("/", routeswitcher.NewProxyHandler())
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
