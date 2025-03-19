package main

import (
	"backend/routes"
	"net/http"
)

func main() {
	r := routes.SetRouter()

	// 8080ポート
	http.ListenAndServe(":8080", r)
}
