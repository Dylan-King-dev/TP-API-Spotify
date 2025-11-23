package main

import (
	"TPSpotify/router"
	"fmt"
	"net/http"
)

func main() {

	router.SetupRoutes()

	assets := http.FileServer(http.Dir("assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	img := http.FileServer(http.Dir("images/"))
	http.Handle("/images/", http.StripPrefix("/images/", img))

	port := ":8080"
	fmt.Println("Server running at http://localhost" + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
	}
}
