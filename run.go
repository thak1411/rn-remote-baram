package main

import (
	"fmt"
	"net/http"

	"rn.com/src/router"
)

func main() {
	mux := router.Get()

	fmt.Println("Running Rn Baram Remote Server")
	http.ListenAndServe(":3003", mux)
}
